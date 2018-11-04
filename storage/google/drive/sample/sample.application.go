package main

import (
	"fmt"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/storage/google/drive"
	"log"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadApplicationConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the S3 extension
	driveExtension := drive.Extension{}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&driveExtension,
	}

	// create and initialize the application
	app := nibbler.Application{}
	if err = app.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	// initialize the drive extension's service separately.
	// this step can be combined with .Init by creating the
	// extension with ConnectServiceOnInit = true
	if err = driveExtension.InitService(); err != nil {
		log.Fatal(err.Error())
	}

	r, err := driveExtension.Srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
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
