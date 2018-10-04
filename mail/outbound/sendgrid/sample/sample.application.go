package extension

import (
	"log"
	"github.com/markdicksonjr/nibbler"
	NibSendGrid "github.com/markdicksonjr/nibbler/mail/outbound/sendgrid"
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

	// allocate the sendgrid extension
	sendgridExtension := NibSendGrid.Extension{}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&sendgridExtension,
	}

	// create and initialize the application
	app := nibbler.Application{}
	err = app.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	response, err := sendgridExtension.SendMail(
		&outbound.Email{Address: "test@example.com", Name: "Example User"},
		"test email",
		&outbound.Email{Address: "mark@example.com", Name: "MD"},
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
