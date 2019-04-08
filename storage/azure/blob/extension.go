package blob

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/markdicksonjr/nibbler"
	"net/url"
)

type Extension struct {
	nibbler.Extension

	accountName string
	accountKey  string
	credential  *azblob.SharedKeyCredential
}

func (s *Extension) Init(app *nibbler.Application) error {
	s.accountName = app.GetConfiguration().Raw.Get("azure", "blob", "account", "name").String("")
	s.accountKey = app.GetConfiguration().Raw.Get("azure", "blob", "account", "key").String("")

	if s.accountName == "" || s.accountKey == "" {
		return errors.New("azure blob extension requires both account name and account key")
	}

	var err error
	s.credential, err = azblob.NewSharedKeyCredential(s.accountName, s.accountKey)

	return err
}

func (s *Extension) AddRoutes(app *nibbler.Application) error {
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return nil
}

func (s *Extension) GetContainerURL(ctx context.Context, name string) (*azblob.ContainerURL, error) {

	// Create a request pipeline that is used to process HTTP(S) requests and responses. It requires
	// your account credentials. In more advanced scenarios, you can configure telemetry, retry policies,
	// logging, and other options. Also, you can configure multiple request pipelines for different scenarios.
	p := azblob.NewPipeline(s.credential, azblob.PipelineOptions{})

	// From the Azure portal, get your Storage account blob service URL endpoint.
	// The URL typically looks like this:
	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", s.accountName))

	if err != nil {
		return nil, err
	}

	// Create an ServiceURL object that wraps the service URL and a request pipeline.
	serviceURL := azblob.NewServiceURL(*u, p)

	// Now, you can use the serviceURL to perform various container and blob operations.

	// This example shows several common operations just to get you started.

	// Create a URL that references a to-be-created container in your Azure Storage account.
	// This returns a ContainerURL object that wraps the container's URL and a request pipeline (inherited from serviceURL)
	containerUrl := serviceURL.NewContainerURL(name)
	return &containerUrl, nil
}
