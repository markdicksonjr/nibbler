package sparkpost

import (
	"errors"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound"
	"github.com/SparkPost/gosparkpost"
)

type Extension struct {
	nibbler.NoOpExtension
	outbound.Sender

	apiKey string
	initialized bool
}

func (s *Extension) Init(app *nibbler.Application) error {
	config := app.GetConfiguration()
	if config == nil || config.Raw == nil {
		return errors.New("sparkpost extension could not get app config")
	}

	configValue := *config.Raw
	s.apiKey = configValue.Get("sparkpost", "api", "key").String("")

	if len(s.apiKey) == 0 {
		return errors.New("sparkpost extension could not get API key")
	}

	s.initialized = true
	return nil
}

func (s *Extension) SendMail(from *outbound.Email, subject string, to []*outbound.Email, plainTextContent string, htmlContent string) (*outbound.Response, error) {
	if !s.initialized {
		return nil, errors.New("sparkpost grid extension used for sending without initialization")
	}

	cfg := &gosparkpost.Config{
		BaseUrl:	"https://api.sparkpost.com",
		ApiKey: 	s.apiKey,
		ApiVersion:	1,
	}
	var sp gosparkpost.Client
	err := sp.Init(cfg)

	if err != nil {
		return nil, err
	}

	content := gosparkpost.Content{
		From:    from.Address, // TODO: apply name
		Subject: "That gopher",
		HTML:    htmlContent,
	}

	var toList []string
	for _, v := range to {
		toList = append(toList, (*v).Address) // TODO: apply name
	}

	tx := &gosparkpost.Transmission{
		Content:    content,
		Recipients: toList,
	}
	_, res, err := sp.Send(tx)

	if res != nil {
		return &outbound.Response{
			Body: "", // TODO: res.Body or res.HTTP.Body via Reader interface
			Headers: res.HTTP.Header,
			StatusCode: res.HTTP.StatusCode,
		}, err
	}

	return nil, err
}
