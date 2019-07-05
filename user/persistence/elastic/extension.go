package elastic

import (
	"context"
	"encoding/json"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/elasticsearch"
)

type Extension struct {
	nibbler.NoOpExtension
	ElasticExtension *elasticsearch.Extension
}

func (s *Extension) Init(app *nibbler.Application) error {
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
			app.GetLogger().Info("in user extension, user index in elastic not acknowledged after creation")
		}
	}

	return nil
}

func (s *Extension) GetUserById(id string) (*nibbler.User, error) {
	ctx := context.Background()
	result, err := s.ElasticExtension.Client.Get().Index("user").Id(id).Do(ctx)

	if err != nil {
		return nil, err
	}

	userValue := nibbler.User{}
	err = json.Unmarshal(*result.Source, &userValue)
	return &userValue, err
}

func (s *Extension) GetUserByEmail(email string) (*nibbler.User, error) {
	ctx := context.Background()

	matchQuery := s.ElasticExtension.NewMatchQuery("email", email)
	result, err := s.ElasticExtension.Client.Search().Index("user").Query(matchQuery).From(0).Size(1).Do(ctx)

	if err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	userValue := nibbler.User{}
	err = json.Unmarshal(*result.Hits.Hits[0].Source, &userValue)
	return &userValue, err
}

func (s *Extension) GetUserByUsername(username string) (*nibbler.User, error) {
	ctx := context.Background()
	matchQuery := s.ElasticExtension.NewMatchQuery("username", username)
	result, err := s.ElasticExtension.Client.Search().Index("user").Query(matchQuery).Size(1).Do(ctx)

	if err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	userValue := nibbler.User{}
	err = json.Unmarshal(*result.Hits.Hits[0].Source, &userValue)
	return &userValue, err
}

func (s *Extension) GetUserByPasswordResetToken(token string) (*nibbler.User, error) {
	ctx := context.Background()
	matchQuery := s.ElasticExtension.NewMatchQuery("passwordResetToken", token)
	result, err := s.ElasticExtension.Client.Search().Index("user").Query(matchQuery).Size(1).Do(ctx)

	if err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	userValue := nibbler.User{}
	err = json.Unmarshal(*result.Hits.Hits[0].Source, &userValue)
	return &userValue, err
}

func (s *Extension) GetUserByEmailValidationToken(token string) (*nibbler.User, error) {
	ctx := context.Background()
	matchQuery := s.ElasticExtension.NewMatchQuery("emailValidationToken", token)
	result, err := s.ElasticExtension.Client.Search().Index("user").Query(matchQuery).Size(1).Do(ctx)

	if err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	userValue := nibbler.User{}
	err = json.Unmarshal(*result.Hits.Hits[0].Source, &userValue)
	return &userValue, err
}

func (s *Extension) Create(userValue *nibbler.User) (*nibbler.User, error) {
	ctx := context.Background()

	// TODO
	_, err := s.ElasticExtension.Client.Index().Index("user").Type("user").Id(userValue.ID).BodyJson(userValue).Do(ctx)
	return nil, err
}

func (s *Extension) Update(userValue *nibbler.User) error {
	ctx := context.Background()

	// TODO
	_, err := s.ElasticExtension.Client.Update().Index("user").Type("user").Id(userValue.ID).Doc(userValue).Do(ctx)

	return err
}

func (s *Extension) UpdatePassword(userValue *nibbler.User) error {
	ctx := context.Background()

	// TODO: password...
	_, err := s.ElasticExtension.Client.Update().Index("user").Type("user").Id(userValue.ID).Doc(userValue).Do(ctx)

	return err
}
