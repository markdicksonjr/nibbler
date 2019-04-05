package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/storage/azure/blob"
	"log"
	"strings"
)

func main() {

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the blob extension
	blobExt := blob.Extension{}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&blobExt,
	}); err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	containerURLPtr, err := blobExt.GetContainerURL(ctx, "testcontainer")
	containerURL := *containerURLPtr

	// create the container on the service (with no metadata and no public access)
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		log.Fatal(err)
	}

	// create a URL that references a to-be-created blob in your Azure Storage account's container.
	// this returns a BlockBlobURL object that wraps the blob's URL and a request pipeline (inherited from containerURL)
	blobURL := containerURL.NewAppendBlobURL("HelloWorld.txt") // Blob names can be mixed case

	_, err = blobURL.Create(ctx, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		log.Fatal(err)
	}

	// create the blob with string (plain text) content.
	data := "Hello World!"
	_, err = blobURL.AppendBlock(ctx, strings.NewReader(data), azblob.AppendBlobAccessConditions{}, nil) //..PutBlob(ctx, strings.NewReader(data), azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		log.Fatal(err)
	}

	// download the blob's contents and verify that it worked correctly
	get, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
		log.Fatal(err)
	}

	downloadedData := &bytes.Buffer{}
	body := get.Body(azblob.RetryReaderOptions{})
	downloadedData.ReadFrom(body)
	body.Close() // The client must close the response body when finished with it
	if data != downloadedData.String() {
		log.Fatal("downloaded data doesn't match uploaded data")
	}

	// list the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
		// get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			log.Fatal(err)
		}

		// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print("Blob name: " + blobInfo.Name + "\n")
		}
	}

	// delete the blob we created earlier.
	_, err = blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		log.Fatal(err)
	}

	// delete the container we created earlier.
	_, err = containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	if err != nil {
		log.Fatal(err)
	}

	// start the app
	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
