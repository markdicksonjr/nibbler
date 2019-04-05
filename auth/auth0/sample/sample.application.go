package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/auth/auth0"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/session/connectors"
	"log"
	"net/http"
)

type SampleExtension struct {
	nibbler.NoOpExtension
	Auth0Extension *auth0.Extension
}

func (s *SampleExtension) AddRoutes(app *nibbler.Application) error {
	app.GetRouter().HandleFunc("/test", s.Auth0Extension.EnforceLoggedIn(s.ProtectedRoute)).Methods("GET")
	return nil
}

func (s *SampleExtension) ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "authorized"}`))
}

func main() {

	// allocate logger and configuration
	config, err := nibbler.LoadConfiguration(nil)

	// allocate session extension
	sessionExtension := session.Extension{
		SessionName: "auth0",
		StoreConnector: connectors.SqlStoreConnector{
			Secret: "something",
		},
	}

	// allocate auth0 extension
	auth0Extension := auth0.Extension{
		SessionExtension:    &sessionExtension,
		LoggedInRedirectUrl: "/",
	}

	// initialize the application, provide config, logger, extensions
	appContext := nibbler.Application{}
	if err := appContext.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&sessionExtension,
		&auth0Extension,
		&SampleExtension{
			Auth0Extension: &auth0Extension,
		},
	}); err != nil {
		log.Fatal(err.Error())
	}

	// run the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
