package s3

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/markdicksonjr/nibbler"
)

type Extension struct {
	nibbler.NoOpExtension

	S3 *s3.S3
}

type S3CredentialProvider struct {
	Config *nibbler.Configuration
}

func (s S3CredentialProvider) Retrieve() (credentials.Value, error) {
	if s.Config == nil || s.Config.Raw == nil {
		return credentials.Value{}, errors.New("app configuration not found by s3 extension")
	}

	return credentials.Value{
		AccessKeyID: s.Config.Raw.Get("s3", "accesskey").String(""),
		SecretAccessKey: s.Config.Raw.Get("s3", "secret").String(""),
		SessionToken:s.Config.Raw.Get("s3", "session", "key").String(""),
		ProviderName: "local",
	}, nil
}

// IsExpired returns if the credentials are no longer valid, and need to be retrieved.
func (s S3CredentialProvider) IsExpired() bool {
	return false
}

func (s *Extension) Init(app *nibbler.Application) error {
	if app.GetConfiguration() == nil || app.GetConfiguration().Raw == nil {
		return errors.New("app configuration not found by s3 extension")
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewCredentials(&S3CredentialProvider{
			Config: app.GetConfiguration(),
		}),
		Endpoint: aws.String(app.GetConfiguration().Raw.Get("s3", "endpoint").String("")),
		Region:   aws.String(app.GetConfiguration().Raw.Get("s3", "region").String("")),
	})

	if err != nil {
		return err
	}

	// create S3 service client
	s.S3 = s3.New(sess)

	return nil
}
