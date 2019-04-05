package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound"
	nibMandrill "github.com/markdicksonjr/nibbler/mail/outbound/mandrill"
	"log"
)

func main() {

	// allocate configuration (from env vars, files, etc)
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the sparkpost extension
	mandrillExtension := nibMandrill.Extension{}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&mandrillExtension,
	}); err != nil {
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
