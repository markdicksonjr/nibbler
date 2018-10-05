package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/auth/auth0"
	"github.com/markdicksonjr/nibbler/session"
	"net/http"
)

type SampleExtension struct {
	Auth0Extension *auth0.Extension
}

func (s *SampleExtension) Init(app *nibbler.Application) error {
	return nil
}

func (s *SampleExtension) AddRoutes(app *nibbler.Application) error {
	app.GetRouter().HandleFunc("/test", s.Auth0Extension.EnforceLoggedIn(s.ProtectedRoute)).Methods("GET")
	return nil
}

func (s *SampleExtension) Destroy(app *nibbler.Application) error {
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
	config, err := nibbler.LoadApplicationConfiguration(nil)

	// allocate session extension
	sessionExtension := session.Extension{
		SessionName: "auth0",
		Secret: "something",
	}

	// allocate auth0 extension
	auth0Extension := auth0.Extension{
		SessionExtension: &sessionExtension,
		LoggedInRedirectUrl: "/",
	}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
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
