package mongo

import (
	"context"
	"github.com/markdicksonjr/nibbler"
	"go.mongodb.org/mongo-driver/mongo"
)

type Extension struct {
	nibbler.NoOpExtension
	Client *mongo.Client
	Url    string
}

// e.g. mongodb://foo:bar@localhost:27017
func (s *Extension) Init(app *nibbler.Application) error {

	// if the Url attribute isn't set, find the config in environment variables
	if len(s.Url) == 0 {
		s.Url = app.GetConfiguration().Raw.Get("mongo", "url").String("")

		if len(s.Url) == 0 {
			s.Url = app.GetConfiguration().Raw.Get("database", "url").String("")
		}
	}

	var err error
	s.Client, err = mongo.NewClient(s.Url)

	if err != nil {
		return err
	}

	return s.Client.Connect(context.TODO())
}
