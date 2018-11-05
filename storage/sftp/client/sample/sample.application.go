package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/storage/sftp/client"
	"log"
	"strconv"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadApplicationConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the SFTP client extension
	sftpExtension := client.Extension{}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&sftpExtension,
	}

	// create and initialize the application
	app := nibbler.Application{}
	if err = app.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	filesInfo, err := sftpExtension.Client.ReadDir("./")
	log.Println(strconv.Itoa(len(filesInfo)))

	// start the app
	err = app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
