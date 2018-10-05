package s3

import (
	"errors"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/markdicksonjr/nibbler"
)

type Extension struct {
	nibbler.Extension

	S3 *s3.S3
}

type S3CredentialProvider struct {
	Config *nibbler.Configuration
}

func (s S3CredentialProvider) Retrieve() (credentials.Value, error) {
	if s.Config == nil || s.Config.Raw == nil {
		return credentials.Value{}, errors.New("app configuration not found by s3 extension")
	}

	configValue := *s.Config.Raw

	return credentials.Value{
		configValue.Get("s3", "accesskey").String(""),
		configValue.Get("s3", "secret").String(""),
		configValue.Get("s3", "session", "key").String(""),
		"local",
	}, nil
}

// IsExpired returns if the credentials are no longer valid, and need to be retrieved.
func (s S3CredentialProvider) IsExpired() bool {
	return false
}

func (s *Extension) Init(app *nibbler.Application) error {
	config := app.GetConfiguration()
	if config == nil || config.Raw == nil {
		return errors.New("app configuration not found by s3 extension")
	}
	configValue := *config.Raw

	creds := credentials.NewCredentials(&S3CredentialProvider{
		Config: config,
	})
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Endpoint: aws.String(configValue.Get("s3", "endpoint").String("")),
		Region: aws.String(configValue.Get("s3", "region").String("")),
	})

	if err != nil {
		log.Fatal(err)
	}

	// create S3 service client
	s.S3 = s3.New(sess)

	return nil
}

func (s *Extension) AddRoutes(app *nibbler.Application) error {
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return nil
}
