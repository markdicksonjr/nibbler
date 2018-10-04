package main

import (
	"net/http"
	"log"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
	"github.com/markdicksonjr/nibbler"
)

type SampleExtension struct {
}

func (s *SampleExtension) Init(context *nibbler.Application) error {
	return nil
}

func (s *SampleExtension) AddRoutes(context *nibbler.Application) error {
	context.GetRouter().HandleFunc("/api/ok", OkResultHandler).Methods("GET")
	return nil
}

func OkResultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "OK"}`))
}

func main() {

	// allocate logger
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	envSources := []source.Source{
		file.NewSource(file.WithPath("./sample.config.json")),
		env.NewSource(),
	}

	// allocate configuration
	config, err := nibbler.LoadApplicationConfiguration(&envSources)

	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&SampleExtension{},
	}

	// initialize the application
	appContext := nibbler.Application{}
	err = appContext.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = appContext.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
