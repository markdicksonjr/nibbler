package main

import (
	"net/http"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/auth/auth0"
	UserAuth0 "github.com/markdicksonjr/nibbler/user/auth/auth0"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/user"
	NibUserSql "github.com/markdicksonjr/nibbler/user/database/sql"
	"log"
)

type SampleExtension struct {
	nibbler.NoOpExtension
	Auth0Extension *UserAuth0.Extension
}

func (s *SampleExtension) AddRoutes(app *nibbler.Application) error {
	app.GetRouter().HandleFunc("/test", s.Auth0Extension.EnforceLoggedIn(s.ProtectedRoute)).Methods("GET")
	return nil
}

func (s *SampleExtension) ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	u, err := s.Auth0Extension.SessionExtension.GetCaller(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "` + err.Error() + `"}`))
		return
	}

	log.Println(u)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "authorized"}`))
}

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}
	config, err := nibbler.LoadConfiguration(nil)

	// allocate session extension
	sessionExtension := session.Extension{
		SessionName: "auth0",
		Secret: "something",
	}

	// prepare models for initialization
	var models []interface{}
	models = append(models, user.User{})

	// allocate the sql extension, with all models
	sqlExtension := sql.Extension{
		Models: models,
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
			SessionExtension: &sessionExtension,
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
	err = appContext.Init(config, &logger, &extensions)

	if err != nil {
		logger.Error(err.Error())
	}

	// start the app
	err = appContext.Run()

	if err != nil {
		logger.Error(err.Error())
	}
}
