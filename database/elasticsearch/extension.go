package elasticsearch

import (
	"errors"
	"github.com/markdicksonjr/nibbler"
	"github.com/olivere/elastic"
)

type Extension struct {
	nibbler.NoOpExtension

	Client *elastic.Client
	Url    string
}

func (s *Extension) Init(app *nibbler.Application) error {
	var err error

	if app.GetConfiguration() == nil {
		return errors.New("app configuration not provided")
	}

	// if the Url attribute isn't set, find the config in environment variables
	if len(s.Url) == 0 {
		s.Url = app.GetConfiguration().Raw.Get("elastic", "url").String("")

		if len(s.Url) == 0 {
			s.Url = app.GetConfiguration().Raw.Get("database", "url").String("http://localhost:9200")
		}
	}

	s.Client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(s.Url))

	return err
}

// TODO: these are pretty silly

func (s *Extension) NewMatchQuery(name string, text interface{}) *elastic.MatchQuery {
	return elastic.NewMatchQuery(name, text)
}

func (s *Extension) NewMatchAllQuery() *elastic.MatchAllQuery {
	return elastic.NewMatchAllQuery()
}

func (s *Extension) NewMatchNoneQuery() *elastic.MatchNoneQuery {
	return elastic.NewMatchNoneQuery()
}

func (s *Extension) NewMatchPhraseQuery(name string, value interface{}) *elastic.MatchPhraseQuery {
	return elastic.NewMatchPhraseQuery(name, value)
}

func (s *Extension) NewBoolQuery() *elastic.BoolQuery {
	return elastic.NewBoolQuery()
}

func (s *Extension) NewBulkDeleteRequest() *elastic.BulkDeleteRequest {
	return elastic.NewBulkDeleteRequest()
}

func (s *Extension) NewBulkIndexRequest() *elastic.BulkIndexRequest {
	return elastic.NewBulkIndexRequest()
}

func (s *Extension) NewBulkUpdateRequest() *elastic.BulkUpdateRequest {
	return elastic.NewBulkUpdateRequest()
}

func (s *Extension) NewIdsQuery(types ...string) *elastic.IdsQuery {
	return elastic.NewIdsQuery()
}
