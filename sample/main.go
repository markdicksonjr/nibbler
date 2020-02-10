package main

import (
	"github.com/markdicksonjr/nibbler"
	"log"
	"net/http"
)

// SampleExtension shows an extremely simple custom extension for nibbler.  It is initially composed of nibbler.NoOpExtension
// in order to allow us to only fill in what we want to for the extension (in this case, only "PostInit").
type SampleExtension struct {
	nibbler.NoOpExtension
}

// PostInit just adds a REST endpoint at "/api/ok" to serve up a simple health-check type message
func (s *SampleExtension) PostInit(context *nibbler.Application) error {
	context.Router.HandleFunc("/api/ok", func(w http.ResponseWriter, _ *http.Request) {
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

	if config.Port == 0 {
		log.Println("WARNING: no port is configured - starting app without http listener")
	}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err := app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&SampleExtension{},
	}); err != nil {
		log.Fatal(err.Error())
	}

	// you could directly interact with your extensions here, if you'd like, as they are all initialized

	// run the app
	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
