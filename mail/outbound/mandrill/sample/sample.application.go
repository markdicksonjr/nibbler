package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound"
	nibMandrill "github.com/markdicksonjr/nibbler/mail/outbound/mandrill"
	"log"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the sparkpost extension
	mandrillExtension := nibMandrill.Extension{}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&mandrillExtension,
	}

	// create and initialize the application
	app := nibbler.Application{}
	if app.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	var toList []*outbound.Email
	toList = append(toList, &outbound.Email{Address: "mark@example.com", Name: "MD"})

	_, err = mandrillExtension.SendMail(
		&outbound.Email{Address: "test@example.com", Name: "Example User"},
		"test email",
		toList,
		"test plain",
		"<strong>test</strong> plain",
	)
}
