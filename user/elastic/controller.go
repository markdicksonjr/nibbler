package sql

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/elasticsearch"
	"github.com/markdicksonjr/nibbler/user"
	"context"
	"strconv"
	"encoding/json"
)

type ElasticPersistenceController struct {
	ElasticExtension *elasticsearch.Extension
}

func (s *ElasticPersistenceController) Init(app *nibbler.Application) error {
	ctx := context.Background()
	exists, err := s.ElasticExtension.Client.IndexExists("user").Do(ctx)

	if err != nil {
		return err
	}

	if !exists {
		createIndex, err := s.ElasticExtension.Client.CreateIndex("user").Do(ctx)

		if err != nil {
			return err
		}

		if !createIndex.Acknowledged {
			(*app.GetLogger()).Info("in user extension, user index in elastic not acknowledged after creation")
		}
	}

	return nil
}

func (s *ElasticPersistenceController) GetUserById(id uint) (*user.User, error) {
	ctx := context.Background()
	result, err := s.ElasticExtension.Client.Get().Index("user").Id(strconv.Itoa(int(id))).Do(ctx)

	if err != nil {
		return nil, err
	}

	userValue := user.User{}
	err = json.Unmarshal(*result.Source, &userValue)
	return &userValue, err
}

func (s *ElasticPersistenceController) GetUserByEmail(email string) (*user.User, error) {
	ctx := context.Background()

	matchQuery := s.ElasticExtension.NewMatchQuery("email", email)
	result, err := s.ElasticExtension.Client.Search().Index("user").Query(matchQuery).From(0).Size(1).Do(ctx)

	if err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	userValue := user.User{}
	err = json.Unmarshal(*result.Hits.Hits[0].Source, &userValue)
	return &userValue, err
}

func (s *ElasticPersistenceController) GetUserByUsername(username string) (*user.User, error) {
	ctx := context.Background()
	matchQuery := s.ElasticExtension.NewMatchQuery("username", username)
	result, err := s.ElasticExtension.Client.Search().Index("user").Query(matchQuery).Size(1).Do(ctx)

	if err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	userValue := user.User{}
	err = json.Unmarshal(*result.Hits.Hits[0].Source, &userValue)
	return &userValue, err
}

func (s *ElasticPersistenceController) GetUserByPasswordResetToken(token string) (*user.User, error) {
	ctx := context.Background()
	matchQuery := s.ElasticExtension.NewMatchQuery("passwordResetExpiration", token)
	result, err := s.ElasticExtension.Client.Search().Index("user").Query(matchQuery).Size(1).Do(ctx)

	if err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	userValue := user.User{}
	err = json.Unmarshal(*result.Hits.Hits[0].Source, &userValue)
	return &userValue, err
}

func (s *ElasticPersistenceController) Create(userValue *user.User) (*user.User, error) {
	ctx := context.Background()

	// TODO
	_, err := s.ElasticExtension.Client.Index().Index("user").Type("user").Id(userValue.ID).BodyJson(userValue).Do(ctx)
	return nil, err
}

func (s *ElasticPersistenceController) Update(userValue *user.User) error {
	ctx := context.Background()

	// TODO
	_, err := s.ElasticExtension.Client.Update().Index("user").Type("user").Id(userValue.ID).Doc(userValue).Do(ctx)

	return err
}

func (s *ElasticPersistenceController) UpdatePassword(userValue *user.User) (error) {
	ctx := context.Background()

	// TODO: password...
	_, err := s.ElasticExtension.Client.Update().Index("user").Type("user").Id(userValue.ID).Doc(userValue).Do(ctx)

	return err
}
