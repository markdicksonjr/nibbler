package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/storage/s3"
	"log"
)

func main() {

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the S3 extension
	s3Extension := s3.Extension{}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&s3Extension,
	}); err != nil {
		log.Fatal(err.Error())
	}

	// list the buckets for the S3 connection we've made
	result, err := s3Extension.S3.ListBuckets(nil)
	if err != nil {
		log.Fatal("Unable to list buckets, " + err.Error())
	}

	// print buckets to the console, with creation time
	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	// start the app
	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
