package main

import (
	"log"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound/sparkpost"
	"github.com/markdicksonjr/nibbler/mail/outbound"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadApplicationConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the sparkpost extension
	sparkpostExtension := sparkpost.Extension{}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&sparkpostExtension,
	}

	// create and initialize the application
	app := nibbler.Application{}
	err = app.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	response, err := sparkpostExtension.SendMail(
		&outbound.Email{Address: "test@example.com", Name: "Example User"},
		"test email",
		&outbound.Email{Address: "markdicksonjr@gmail.com", Name: "MD"},
		"test plain",
		"<strong>test</strong> plain",
	)

	log.Println(response)

	if err != nil {
		log.Fatal(err.Error())
	}

	// start the app
	err = app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
