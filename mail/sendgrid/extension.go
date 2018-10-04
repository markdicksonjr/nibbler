package sendgrid

import (
	"errors"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/markdicksonjr/nibbler"
)

type Extension struct {
	apiKey string
	initialized bool
}

func (s *Extension) Init(app *nibbler.Application) error {
	config := app.GetConfiguration()
	if config == nil || config.Raw == nil {
		return errors.New("sendgrid extension could not get app config")
	}

	configValue := *config.Raw
	s.apiKey = configValue.Get("sendgrid", "api", "key").String("")

	if len(s.apiKey) == 0 {
		return errors.New("sendgrid extension could not get API key")
	}

	s.initialized = true
	return nil
}

func (s *Extension) AddRoutes(app *nibbler.Application) error {
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return nil
}

func (s *Extension) SendMail(from *mail.Email, subject string, to *mail.Email, plainTextContent string, htmlContent string) (*rest.Response, error) {
	if !s.initialized {
		return nil, errors.New("send grid extension used for sending without initialization")
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)
	return client.Send(message)
}