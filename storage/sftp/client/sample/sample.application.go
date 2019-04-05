package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/storage/sftp/client"
	"log"
	"strconv"
)

func main() {

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the SFTP client extension
	sftpExtension := client.Extension{}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&sftpExtension,
	}); err != nil {
		log.Fatal(err.Error())
	}

	filesInfo, err := sftpExtension.Client.ReadDir("./")
	log.Println(strconv.Itoa(len(filesInfo)))

	// start the app
	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
