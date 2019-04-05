package main

import (
	"fmt"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/storage/google/drive"
	"log"
)

func main() {

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the S3 extension
	driveExtension := drive.Extension{}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&driveExtension,
	}); err != nil {
		log.Fatal(err.Error())
	}

	// initialize the drive extension's service separately.
	// this step can be combined with .Init by creating the
	// extension with ConnectServiceOnInit = true
	if err = driveExtension.InitService(); err != nil {
		log.Fatal(err.Error())
	}

	r, err := driveExtension.Srv.Files.List().PageSize(25).
		Fields("nextPageToken, files(id, name, md5Checksum, mimeType, parents)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	fmt.Println("Files/Folders:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			if i.MimeType == "application/vnd.google-apps.folder" {
				fmt.Printf("Folder: %s [%s] (%s)\n", i.Name, i.MimeType, i.Id)
			} else {
				fmt.Printf("File: %s [%s] (%s) md5:%s\n", i.Name, i.MimeType, i.Id, i.Md5Checksum)
			}
		}
	}
	//
	//// print buckets to the console, with creation time
	//fmt.Println("Buckets:")
	//for _, b := range result.Buckets {
	//	fmt.Printf("* %s created on %s\n",
	//		aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	//}
	//
	//// start the app
	//err = app.Run()
	//
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
}
