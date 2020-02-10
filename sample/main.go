package main

import (
	"github.com/markdicksonjr/nibbler"
	"net/http"
)

// SampleExtension shows an extremely simple custom extension for nibbler.  It is initially composed of
// nibbler.NoOpExtension in order to allow us to only fill in what we want to for the extension (in this case,
// only "PostInit").
//
// The extension serves up static content at "/" and a simple API endpoint at "/api/ok"
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

	// allocate logger
	logger := nibbler.DefaultLogger{}

	// allocate configuration
	config, err := nibbler.LoadConfiguration()

	// use a nibbler utility function to handle the error - if it's non-nil, it will log at error level and exit
	nibbler.LogFatalNonNil(logger, err, "while loading configuration")

	// display a notice if the configuration doesn't specify a port - refer to README to learn about configuration
	// note that nibbler doesn't NEED a port to run, but our example sets up a route, which won't be reachable
	// without specifying a port.  As such, this isn't typically something a nibbler app will need to check
	if config.Port == 0 {
		logger.Warn("no port is configured - starting app without http listener")
	}

	// allocate the application
	app := nibbler.Application{}

	// initialize the application, provide config, logger, extensions - any error is fatal
	nibbler.LogFatalNonNil(logger, app.Init(config, logger, []nibbler.Extension{
		&SampleExtension{},
	}))

	// you could directly interact with your extensions here, if you'd like, as they are all initialized

	// run the app, any error is fatal - this will block until app destruction
	nibbler.LogFatalNonNil(logger, app.Run())
}
