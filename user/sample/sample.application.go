package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/elasticsearch"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/user"
	NibUserElastic "github.com/markdicksonjr/nibbler/user/persistence/elastic"
	NibUserSql "github.com/markdicksonjr/nibbler/user/persistence/sql"
	"log"
)

type UserAndDbExtensions struct {
	UserExtension            *user.Extension
	UserPersistenceExtension user.PersistenceExtension
	DbExtension              nibbler.Extension
}

func allocateSqlExtensions() UserAndDbExtensions {

	// allocate an SQL extension, providing the models for auto-migration
	sqlExtension := sql.Extension{Models: []interface{}{
		nibbler.User{},
	}}

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
		DbExtension:              &elasticExtension,
		UserPersistenceExtension: &elasticUserExtension,
		UserExtension: &user.Extension{
			PersistenceExtension: &elasticUserExtension,
		},
	}
}

func main() {

	// allocate logger and configuration
	config, err := nibbler.LoadConfiguration(nil)

	sqlExtensions := allocateSqlExtensions()

	// initialize the application, provide config, logger, extensions
	appContext := nibbler.Application{}
	if err = appContext.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		sqlExtensions.DbExtension,
		sqlExtensions.UserPersistenceExtension,
		sqlExtensions.UserExtension,
	}); err != nil {
		log.Fatal(err.Error())
	}

	// create a test user
	emailVal := "test@example.com"
	_, errCreate := sqlExtensions.UserExtension.Create(&nibbler.User{
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

	if err = sqlExtensions.UserExtension.Update(uV); err != nil {
		log.Fatal(err.Error())
	}

	// start the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
