package nibbler

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/config"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Extension interface {
	Init(app *Application) error
	AddRoutes(app *Application) error
	Destroy(app *Application) error
}

type Logger interface {
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
	config     *Configuration
	extensions []Extension
	logger     Logger
	router     *mux.Router
	stopSignal chan os.Signal
	mode       string
}

func (ac *Application) Init(config *Configuration, logger Logger, extensions []Extension) error {
	ac.config = config
	ac.logger = logger
	ac.extensions = extensions
	ac.mode = ac.config.Raw.Get("nibbler", "mode").String("web")

	// prepare a general-use error variable
	var err error

	// initialize all extensions
	// if any error occurred, return the error and stop processing
	for _, x := range extensions {
		if err = x.Init(ac); err != nil {
			return err
		}
	}

	// get the configured app mode, accounting for the default value
	// only add routes if it's a web mode app
	if ac.mode == "web" {

		// allocate a router
		ac.router = mux.NewRouter()

		// init extension routes
		// if any error occurred, return the error and stop processing
		for _, x := range extensions {
			if err = x.AddRoutes(ac); err != nil {
				return err
			}
		}

		// set up the static directory routing
		ac.router.PathPrefix("/").Handler(http.FileServer(http.Dir(config.StaticDirectory)))
		http.Handle("/", ac.router)
	}

	return nil
}

func (ac *Application) Run() error {
	var err error

	// allocate and prep signal channel (listen for some stop signals from the OS)
	ac.stopSignal = make(chan os.Signal, 1)
	signal.Notify(ac.stopSignal, syscall.SIGINT, syscall.SIGTERM)

	if ac.mode == "web" {

		// allocate a server
		h := &http.Server{Addr: ":" + strconv.Itoa(ac.config.Port), Handler: nil}

		// run our server listener in a goroutine
		go func() {
			ac.logger.Info("starting server")
			err = startServer(h, ac)
		}()

		// wait for a signal
		<-ac.stopSignal

		// shut down the server
		ac.logger.Info("shutting down the server")

		// log a shutdown error (factor into return value later)
		if err = h.Shutdown(context.Background()); err != nil {
			ac.logger.Error(err.Error())
		}
	} else if ac.mode == "worker" {

		// wait for a signal
		<-ac.stopSignal

	} else {
		return errors.New("unknown nibbler mode detected: " + ac.mode)
	}

	ac.logger.Info("shutting down the application")

	// destroy extensions in reverse order
	for i := range ac.extensions {
		err = ac.extensions[len(ac.extensions)-i-1].Destroy(ac)
	}

	// return any (the latest, which could be an improvement we could make) extension destroy error
	return err
}

func startServer(h *http.Server, ac *Application) error {

	// log that we're listening and state the port
	ac.logger.Info("listening on " + strconv.Itoa((*ac.config).Port))

	// listen (this blocks) - log an error if it happened
	if err := h.ListenAndServe(); err != nil {
		ac.logger.Error("failed to initialize server: " + err.Error())
		return err
	}

	return nil
}

func (ac *Application) GetLogger() Logger {
	return ac.logger
}

func (ac *Application) GetConfiguration() *Configuration {
	return ac.config
}

func (ac *Application) GetRouter() *mux.Router {
	return ac.router
}
