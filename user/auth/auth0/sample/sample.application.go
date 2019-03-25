package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/auth/auth0"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/session/connectors"
	"github.com/markdicksonjr/nibbler/user"
	UserAuth0 "github.com/markdicksonjr/nibbler/user/auth/auth0"
	NibUserSql "github.com/markdicksonjr/nibbler/user/database/sql"
	"log"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}
	config, err := nibbler.LoadConfiguration(nil)
	if err != nil {
		log.Fatal(err)
	}

	// allocate the sql extension, with all models
	// if database.uri is not configured, an in-memory SQLite database will be used
	sqlExtension := sql.Extension{
		Models: []interface{}{
			user.User{},
		},
	}

	// allocate session extension
	sessionExtension := session.Extension{
		SessionName: "auth0",
		StoreConnector: connectors.SqlStoreConnector{
			SqlExtension: &sqlExtension,
			Secret:       "somesecret",
		},
	}

	// allocate user extension, providing sql extension to it
	userExtension := user.Extension{
		PersistenceExtension: &NibUserSql.Extension{
			SqlExtension: &sqlExtension,
		},
	}

	// allocate user auth0 extension
	auth0Extension := UserAuth0.Extension{
		Extension: auth0.Extension{
			SessionExtension:    &sessionExtension,
			LoggedInRedirectUrl: "/",
		},
		UserExtension: &userExtension,
	}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&sqlExtension,
		&userExtension,
		&sessionExtension,
		&auth0Extension,
		&SampleExtension{
			Auth0Extension: &auth0Extension,
		},
	}

	// initialize the application
	appContext := nibbler.Application{}
	if err = appContext.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	// create a test user, if it does not exist
	emailVal := "someone@example.com"
	password := ""
	_, errCreate := userExtension.Create(&user.User{
		Email:    &emailVal,
		Password: &password,
	})

	// assert the test user got created
	if errCreate != nil {
		log.Fatal(errCreate.Error())
	}

	// start the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
