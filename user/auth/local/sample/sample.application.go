package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/session/connectors"
	NibUser "github.com/markdicksonjr/nibbler/user"
	NibUserLocalAuth "github.com/markdicksonjr/nibbler/user/auth/local"
	NibUserSql "github.com/markdicksonjr/nibbler/user/database/sql"
	"log"
	"net/http"
)

type SampleExtension struct {
	nibbler.NoOpExtension
	AuthExtension *NibUserLocalAuth.Extension
}

func (s *SampleExtension) AddRoutes(app *nibbler.Application) error {
	app.GetRouter().HandleFunc("/api/test", s.AuthExtension.EnforceLoggedIn(s.ProtectedRoute)).Methods("POST")
	return nil
}

func (s *SampleExtension) ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	nibbler.Write200Json(w, `{"result": "authorized"}`)
}

func main() {

	// allocate configuration
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err)
	}

	// allocate the sql extension, with all models
	sqlExtension := sql.Extension{
		Models: []interface{}{
			NibUser.User{},
		},
	}

	// allocate user extension, providing sql extension to it
	userExtension := NibUser.Extension{
		PersistenceExtension: &NibUserSql.Extension{
			SqlExtension: &sqlExtension,
		},
	}

	// allocate session extension
	connector := connectors.SqlStoreConnector{
		SqlExtension: &sqlExtension,
		Secret:       "dumbsecret",
	}
	sessionExtension := session.Extension{
		StoreConnector: &connector,
		SessionName:    "dumbcookie",
	}

	// allocate the sendgrid extension
	//sendgridExtension := sendgrid.Extension{}

	// allocate user local auth extension
	userLocalAuthExtension := NibUserLocalAuth.Extension{
		SessionExtension: &sessionExtension,
		UserExtension:    &userExtension,
		//Sender:                 &sendgridExtension,
		PasswordResetEnabled:   false,
		PasswordResetFromName:  "Test",
		PasswordResetFromEmail: "test@example.com",
		PasswordResetRedirect:  "http://localhost:3000/test-ui",

		RegistrationEnabled:        true,
		EmailVerificationEnabled:   false,
		EmailVerificationFromName:  "Test",
		EmailVerificationFromEmail: "test@example.com",
		EmailVerificationRedirect:  "http://localhost:3000/verify",
	}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&sqlExtension,
		&userExtension,
		&sessionExtension,
		&userLocalAuthExtension,
		//&sendgridExtension,
		&SampleExtension{
			AuthExtension: &userLocalAuthExtension,
		},
	}); err != nil {
		log.Fatal(err.Error())
	}

	// check to see if our test user exists
	emailVal := "markdicksonjr@gmail.com"
	user, errGet := userExtension.GetUserByEmail(emailVal)

	if errGet != nil {
		log.Fatal(errGet)
	}

	if user == nil {

		// create a test user, if it does not exist
		password, err := NibUserLocalAuth.GeneratePasswordHash("tester123")

		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = userExtension.Create(&NibUser.User{
			Email:    &emailVal,
			Password: &password,
		})

		// assert the test user got created
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// start the app
	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
