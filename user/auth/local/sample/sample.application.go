package main

import (
	"log"
	"net/http"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound/sendgrid"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/database/sql"
	NibUser "github.com/markdicksonjr/nibbler/user"
	NibUserLocalAuth "github.com/markdicksonjr/nibbler/user/auth/local"
	NibUserSql "github.com/markdicksonjr/nibbler/user/database/sql"
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "authorized"}`))
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
	models = append(models, NibUser.User{})

	// allocate the sql extension, with all models
	sqlExtension := sql.Extension{
		Models: models,
	}

	// allocate user extension, providing sql extension to it
	userExtension := NibUser.Extension{
		PersistenceController: &NibUserSql.SqlExtension{
			SqlExtension: &sqlExtension,
		},
	}

	// allocate session extension
	sessionExtension := session.Extension{
		Secret:      "dumbsecret",
		SessionName: "dumbcookie",
	}

	// allocate the sendgrid extension
	sendgridExtension := sendgrid.Extension{}

	// allocate user local auth extension
	userLocalAuthExtension := NibUserLocalAuth.Extension{
		SessionExtension:       &sessionExtension,
		UserExtension:          &userExtension,
		Sender:     			&sendgridExtension,
		PasswordResetEnabled:   true,
		PasswordResetFromName:  "Test",
		PasswordResetFromEmail: "test@example.com",
		PasswordResetRedirect:  "http://localhost:3000/test-ui",
	}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&sqlExtension,
		&userExtension,
		&sessionExtension,
		&userLocalAuthExtension,
		&sendgridExtension,
		&SampleExtension{
			AuthExtension: &userLocalAuthExtension,
		},
	}

	// initialize the application
	app := nibbler.Application{}
	err = app.Init(config, &logger, &extensions)

	if err != nil {
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
		password := NibUserLocalAuth.GeneratePasswordHash("tester123")
		_, errCreate := userExtension.Create(&NibUser.User{
			Email: &emailVal,
			Password: &password,
		})

		// assert the test user got created
		if errCreate != nil {
			log.Fatal(errCreate.Error())
		}
	}

	// start the app
	err = app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
