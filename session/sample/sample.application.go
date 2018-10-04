package main

import (
	"log"
	"github.com/gorilla/sessions"
	"github.com/wader/gormstore"
	"github.com/jinzhu/gorm"
	_ "github.com/michaeljs1990/sqlitestore"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/sql"
	"github.com/markdicksonjr/nibbler/user"
	NibUserSql "github.com/markdicksonjr/nibbler/user/sql"
)

type SqlMemoryStoreConnector struct {
}

func (s SqlMemoryStoreConnector) Connect() (error, sessions.Store) {
	db, err := gorm.Open("sqlite3", ":memory:")

	if err != nil {
		return err, nil
	}

	store := gormstore.NewOptions(db,
		gormstore.Options{},
		[]byte("some-key"),
	)

	return nil, store
}

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration
	config, err := nibbler.LoadApplicationConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// prepare models for initialization
	var models []interface{}
	models = append(models, user.User{})

	// allocate an SQL controller, providing an sql extension
	sqlController := NibUserSql.SqlPersistenceController{
		SqlExtension: &sql.Extension{
			Models: models,
		},
	}

	// allocate user extension, providing sql extension to it
	userExtension := user.Extension{
		PersistenceController:  &sqlController, // &elasticController,
	}

	// allocate session extension, with an optional custom connector
	var sessionConnector session.SessionStoreConnector = &SqlMemoryStoreConnector{}
	sessionExtension := session.Extension{
		StoreConnector: &sessionConnector,
		Secret: "something",
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
