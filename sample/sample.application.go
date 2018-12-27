package main

import (
	"github.com/markdicksonjr/nibbler"
	"log"
	"net/http"
)

type SampleExtension struct {
	nibbler.NoOpExtension
}

func (s *SampleExtension) AddRoutes(context *nibbler.Application) error {
	context.GetRouter().HandleFunc("/api/ok", func (w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "OK"}`))
	}).Methods("GET")
	return nil
}

func main() {

	// allocate logger
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration
	config, err := nibbler.LoadConfiguration(nil)

	// any error is fatal at this point
	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&SampleExtension{},
	}

	// initialize the application
	appContext := nibbler.Application{}
	if err := appContext.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	// run the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
