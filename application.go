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

	// Init should handle the initialization of the extension with the knowledge that the Router and other Extensions
	// are not likely to be initialized yet
	Init(app *Application) error

	// PostInit should handle all initialization that could not be handled in Init.  Specifically, anything requiring
	// the Router and other Extensions to be initialized
	PostInit(app *Application) error

	// Destroy handles application shutdown, with the most recently-initialized extensions destroyed first
	Destroy(app *Application) error
}

type Logger interface {
	Trace(message ...string)
	Debug(message ...string)
	Error(message ...string)
	Info(message ...string)
	Warn(message ...string)
}

type Configuration struct {
	HeaderConfiguration HeaderConfiguration
	Port                int
	Raw                 config.Config
	StaticDirectory     string
}

type HeaderConfiguration struct {
	AccessControlAllowHeaders string
	AccessControlAllowMethods string
	AccessControlAllowOrigin  string
}

type Application struct {
	Config     *Configuration
	Router     *mux.Router
	Logger     Logger
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
		if err = x.Init(ac); err != nil {
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
			if err = x.PostInit(ac); err != nil {
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
			if err = x.PostInit(ac); err != nil {
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
			ac.Logger.Info("starting server")
			err = startServer(h, ac)
		}()

		// wait for a signal
		<-ac.stopSignal

		// shut down the server
		ac.Logger.Info("shutting down the server")

		// log a shutdown error (factor into return value later)
		if err = h.Shutdown(context.Background()); err != nil {
			ac.Logger.Error(err.Error())
		}
	} else {

		// wait for a signal
		<-ac.stopSignal

	}

	ac.Logger.Info("shutting down the application")

	// destroy extensions in reverse order
	for i := range ac.extensions {
		err = ac.extensions[len(ac.extensions)-i-1].Destroy(ac)
	}

	// return any (the latest, which could be an improvement we could make) extension destroy error
	return err
}

func startServer(h *http.Server, ac *Application) error {

	// log that we're listening and state the port
	ac.Logger.Info("listening on " + strconv.Itoa((*ac.Config).Port))

	// listen (this blocks) - log an error if it happened
	if err := h.ListenAndServe(); err != nil {
		ac.Logger.Error("failed to initialize server: " + err.Error())
		return err
	}

	return nil
}
