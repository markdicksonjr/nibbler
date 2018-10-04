package extension

import (
	"log"
	"github.com/markdicksonjr/nibble"
	NibSendGrid "github.com/markdicksonjr/nibble-sendgrid/extension"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {

	// allocate logger and configuration
	var logger nibble.Logger = nibble.DefaultLogger{}

	// allocate configuration (from env vars, files, etc)
	config, err := nibble.LoadApplicationConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the sendgrid extension
	sendgridExtension := NibSendGrid.Extension{}

	// prepare extensions for initialization
	extensions := []nibble.Extension{
		&sendgridExtension,
	}

	// create and initialize the application
	app := nibble.Application{}
	err = app.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	response, err := sendgridExtension.SendMail(
		mail.NewEmail("Example User", "test@example.com"),
		"test email",
		mail.NewEmail("MD", "markdicksonjr@gmail.com"),
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
