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
	context.GetRouter().HandleFunc("/api/ok", func(w http.ResponseWriter, _ *http.Request) {
		nibbler.Write200Json(w, `{"result": "OK"}`)
	}).Methods("GET")
	return nil
}

func main() {

	// allocate configuration
	config, err := nibbler.LoadConfiguration()

	// any error is fatal at this point
	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize the application, provide config, logger, extensions
	appContext := nibbler.Application{}
	if err := appContext.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&SampleExtension{},
	}); err != nil {
		log.Fatal(err.Error())
	}

	// run the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
