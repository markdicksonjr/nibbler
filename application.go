package nibbler

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/config"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Extension interface {

	// GetName returns a static name for the extension, primarily for logging purposes
	GetName() string

	// Init should handle the initialization of the extension with the knowledge that the Router and other Extensions
	// are not likely to be initialized yet
	Init(app *Application) error

	// PostInit should handle all initialization that could not be handled in Init.  Specifically, anything requiring
	// the Router and other Extensions to be initialized
	PostInit(app *Application) error

	// Destroy handles application shutdown, with the most recently-initialized extensions destroyed first
	Destroy(app *Application) error
}

// Logger is a generic interface to reflect logging output at various levels
type Logger interface {
	Trace(message ...string)
	Debug(message ...string)
	Error(message ...string)
	Info(message ...string)
	Warn(message ...string)
}

type Configuration struct {
	Headers         HeaderConfiguration
	Port            int
	Raw             config.Config
	StaticDirectory string
}

type HeaderConfiguration struct {
	AccessControlAllowHeaders string
	AccessControlAllowMethods string
	AccessControlAllowOrigin  string
}

type Application struct {
	Config     *Configuration
	Logger     Logger
	Router     *mux.Router
	extensions []Extension
	stopSignal chan os.Signal
}

func (ac *Application) Init(config *Configuration, logger Logger, extensions []Extension) error {
	ac.Config = config
	ac.Logger = logger
	ac.extensions = extensions

	// prepare a general-use error variable
	var err error

	// initialize all extensions
	for _, x := range extensions {

		// if any error occurred, return the error and stop processing
		if err = LogErrorNonNil(logger, x.Init(ac)); err != nil {
			return err
		}
	}

	// if a port is configured, set up our listener and init extensions
	if ac.Config.Port != 0 {

		// allocate a router
		ac.Router = mux.NewRouter()

		// init extension routes
		for _, x := range extensions {

			// if any error occurred, return the error and stop processing
			if err = LogErrorNonNil(logger, x.PostInit(ac), "while running Init on extension \"" + x.GetName() + "\""); err != nil {
				return err
			}
		}

		// set up the static directory routing
		ac.Router.PathPrefix("/").Handler(http.FileServer(http.Dir(config.StaticDirectory)))
		http.Handle("/", ac.Router)
	} else {

		// init extension routes
		for _, x := range extensions {

			// if any error occurs, return the error and stop processing
			if err = LogErrorNonNil(ac.Logger, x.PostInit(ac), "while running PostInit on extension \"" + x.GetName() + "\""); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ac *Application) Run() error {
	var err error

	// allocate and prep signal channel (listen for some stop signals from the OS)
	ac.stopSignal = make(chan os.Signal, 1)
	signal.Notify(ac.stopSignal, syscall.SIGINT, syscall.SIGTERM)

	if ac.Config.Port != 0 {

		// allocate a server
		h := &http.Server{Addr: ":" + strconv.Itoa(ac.Config.Port), Handler: nil}

		// run our server listener in a goroutine
		go func() {

			// log that we're listening and state the port
			ac.Logger.Info("listening on " + strconv.Itoa((*ac.Config).Port))

			// listen (this blocks) - log an error if it happened
			LogFatalNonNil(ac.Logger, h.ListenAndServe(), "failed to initialize server")
		}()

		// wait for a signal
		<-ac.stopSignal

		// shut down the server
		ac.Logger.Info("shutting down")

		// log a shutdown error (factor into return value later)
		if err = LogErrorNonNil(ac.Logger, h.Shutdown(context.Background()), "while shutting down"); err != nil {
			ac.Logger.Error(err.Error())
		}
	} else {

		// wait for a signal
		<-ac.stopSignal
	}

	ac.Logger.Info("shutting down the application")

	// destroy extensions in reverse order (keep going on error to try to close as much as we can)
	for i := range ac.extensions {
		x := ac.extensions[len(ac.extensions)-i-1]
		destroyErr := LogErrorNonNil(ac.Logger, x.Destroy(ac), "while destroying extension \"" + x.GetName() + "\"")
		if destroyErr != nil {
			err = destroyErr
		}
	}

	// return any (the latest, which could be an improvement we could make) extension destroy error
	return err
}
