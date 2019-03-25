package mandrill

import (
	"github.com/keighl/mandrill"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound"
	"github.com/pkg/errors"
	"log"
)

type Extension struct {
	nibbler.NoOpExtension
	outbound.Sender

	initialized bool

	Key string

	Client *mandrill.Client
}

func (s *Extension) Init(app *nibbler.Application) error {
	if len(s.Key) == 0 {
		s.Key = (*app.GetConfiguration().Raw).Get("mandrill", "api", "key").String("")
	}

	if len(s.Key) == 0 {
		return errors.New("mandrill.api.key must be defined for nibbler mandrill extension")
	}

	err := s.Connect()
	s.initialized = true
	return err
}

func (s *Extension) Connect() error {
	s.Client = mandrill.ClientWithKey(s.Key)

	return nil
}

func (s *Extension) SendMail(from *outbound.Email, subject string, to []*outbound.Email, plainTextContent string, htmlContent string) (*outbound.Response, error) {
	if !s.initialized {
		return nil, errors.New("send grid extension used for sending without initialization")
	}

	if from == nil || len((*from).Address) == 0 {
		return nil, errors.New("send grid extension requires 'from' field")
	}

	if len(to) == 0 || to[0] == nil || len((*to[0]).Address) == 0 {
		return nil, errors.New("send grid extension requires at least one recipient")
	}

	message := &mandrill.Message{}
	message.AddRecipient((*to[0]).Address, (*to[0]).Name, "to")
	message.FromEmail = (*from).Address
	message.FromName = (*from).Name
	message.Subject = subject
	message.HTML = htmlContent
	message.Text = plainTextContent

	if _, err := s.Client.MessagesSend(message); err != nil {
		log.Fatal(err.Error())
	}

	return nil, nil
}
