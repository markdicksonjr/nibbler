package main

import (
	"log"
	"github.com/markdicksonjr/nibbler/user"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/database/elasticsearch"
	NibUserElastic "github.com/markdicksonjr/nibbler/user/database/elastic"
	NibUserSql "github.com/markdicksonjr/nibbler/user/database/sql"
)

type UserAndDbExtensions struct {
	UserExtension 				*user.Extension
	UserPersistenceExtension 	user.PersistenceExtension
	DbExtension 				nibbler.Extension
}

func allocateSqlExtensions() UserAndDbExtensions {

	// prepare models for initialization
	var models []interface{}
	models = append(models, user.User{})

	// allocate an SQL extension, providing the models for auto-migration
	sqlExtension := sql.Extension{ Models: models }

	sqlUserExtension := NibUserSql.Extension{
		SqlExtension: &sqlExtension,
	}

	return UserAndDbExtensions{
		DbExtension: &sqlExtension,
		UserExtension: &user.Extension{
			PersistenceExtension: &sqlUserExtension,
		},
		UserPersistenceExtension: &sqlUserExtension,
	}
}

func allocateEsExtensions() UserAndDbExtensions {
	elasticExtension := elasticsearch.Extension{}

	elasticUserExtension := NibUserElastic.Extension{
		ElasticExtension: &elasticExtension,
	}

	return UserAndDbExtensions{
		DbExtension: &elasticExtension,
		UserPersistenceExtension: &elasticUserExtension,
		UserExtension: &user.Extension{
			PersistenceExtension: &elasticUserExtension,
		},
	}
}

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}
	config, err := nibbler.LoadApplicationConfiguration(nil)

	sqlExtensions := allocateSqlExtensions()

	// prepare extension(s) for initialization
	extensions := []nibbler.Extension{
		sqlExtensions.DbExtension,
		sqlExtensions.UserPersistenceExtension,
		sqlExtensions.UserExtension,
	}

	// initialize the application context
	appContext := nibbler.Application{}
	err = appContext.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	// create a test user
	emailVal := "test@example.com"
	_, errCreate := sqlExtensions.UserExtension.Create(&user.User{
		Email: &emailVal,
	})

	if errCreate != nil {
		log.Fatal(errCreate.Error())
	}

	uV, errFind := sqlExtensions.UserExtension.GetUserByEmail(emailVal)

	if errFind != nil {
		log.Fatal(errFind.Error())
	}

	log.Println(uV)

	firstName := "testfirst"
	lastName := "testlast"
	uV.FirstName = &firstName
	uV.LastName = &lastName
	err = sqlExtensions.UserExtension.Update(uV)

	if err != nil {
		log.Fatal(err.Error())
	}

	// start the app
	err = appContext.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
