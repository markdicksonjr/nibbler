package sendgrid

import (
	"errors"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound"
)

type Extension struct {
	outbound.Sender
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

func (s *Extension) SendMail(from *outbound.Email, subject string, to *outbound.Email, plainTextContent string, htmlContent string) (*outbound.Response, error) {
	if !s.initialized {
		return nil, errors.New("send grid extension used for sending without initialization")
	}

	fromSg := mail.Email{
		Name: (*from).Name,
		Address: (*from).Address,
	}

	toSg := mail.Email{
		Name: (*to).Name,
		Address: (*to).Address,
	}

	message := mail.NewSingleEmail(&fromSg, subject, &toSg, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)
	res, err := client.Send(message)

	if res != nil {
		return &outbound.Response{
			Body: res.Body,
			Headers: res.Headers,
			StatusCode: res.StatusCode,
		}, err
	}

	return nil, err
}