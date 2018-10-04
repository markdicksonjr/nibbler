package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"github.com/Azure/azure-storage-blob-go/2016-05-31/azblob"
	"github.com/markdicksonjr/nibbler/storage/azure/blob"
	"github.com/markdicksonjr/nibbler"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadApplicationConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the blob extension
	blobExt := blob.Extension{}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&blobExt,
	}

	// create and initialize the application
	app := nibbler.Application{}
	err = app.Init(config, &logger, &extensions)

	if err != nil {
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
	blobURL := containerURL.NewBlockBlobURL("HelloWorld.txt") // Blob names can be mixed case

	// create the blob with string (plain text) content.
	data := "Hello World!"
	_, err = blobURL.PutBlob(ctx, strings.NewReader(data), azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		log.Fatal(err)
	}

	// download the blob's contents and verify that it worked correctly
	get, err := blobURL.GetBlob(ctx, azblob.BlobRange{}, azblob.BlobAccessConditions{}, false)
	if err != nil {
		log.Fatal(err)
	}

	downloadedData := &bytes.Buffer{}
	downloadedData.ReadFrom(get.Body())
	get.Body().Close() // The client must close the response body when finished with it
	if data != downloadedData.String() {
		log.Fatal("downloaded data doesn't match uploaded data")
	}

	// list the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
		// get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobs(ctx, marker, azblob.ListBlobsOptions{})
		if err != nil {
			log.Fatal(err)
		}

		// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Blobs.Blob {
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
	err = app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}


