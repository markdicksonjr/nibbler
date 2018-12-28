package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/session/connectors"
	"github.com/markdicksonjr/nibbler/user"
	NibUserSql "github.com/markdicksonjr/nibbler/user/database/sql"
	_ "github.com/michaeljs1990/sqlitestore"
	"log"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// prepare models for initialization
	var models []interface{}
	models = append(models, user.User{})

	// allocate an SQL controller, providing an sql extension
	sqlController := NibUserSql.Extension{
		SqlExtension: &sql.Extension{
			Models: models,
		},
	}

	// allocate user extension, providing sql extension to it
	userExtension := user.Extension{
		PersistenceExtension:  &sqlController, // &elasticController,
	}

	// allocate session extension, with an optional custom connector
	var sessionConnector session.SessionStoreConnector = &connectors.SqlMemoryStoreConnector{
		Secret: "somesecret",
		SqlExtension: sqlController.SqlExtension,
	}
	sessionExtension := session.Extension{
		StoreConnector: sessionConnector,
	}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		sqlController.SqlExtension,
		&userExtension,
		&sessionExtension,
	}

	// initialize the application
	appContext := nibbler.Application{}
	err = appContext.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	// start the app
	err = appContext.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
