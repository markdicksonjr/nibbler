package nibbler

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/micro/go-config"
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
	Debug(message string)
	Error(message string)
	Info(message string)
	Warn(message string)
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
}

func (ac *Application) Init(config *Configuration, logger Logger, extensions []Extension) error {
	ac.config = config
	ac.logger = logger
	ac.extensions = extensions

	// dereference parameters for ease-of-use
	configValue := *config

	// prepare a general-use error variable
	var err error = nil

	// initialize all extensions
	// if any error occurred, return the error and stop processing
	for _, x := range extensions {
		if err = x.Init(ac); err != nil {
			return err
		}
	}

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
	ac.router.PathPrefix("/").Handler(http.FileServer(http.Dir(configValue.StaticDirectory)))

	logger.Info("Starting server")

	http.Handle("/", ac.router)

	return nil
}

func (ac *Application) Run() error {

	// get the configured app mode, accounting for the default value
	mode := ac.config.Raw.Get("nibbler", "mode").String("web")

	// allocate and prep signal channel
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// prepare a variable to be used for something going wrong during shutdown, regardless of mode
	var shutdownError error

	if mode == "web" {

		// allocate a server
		h := &http.Server{Addr: ":" + strconv.Itoa(ac.config.Port), Handler: nil}

		// run our server listener in a goroutine
		go startServer(h, ac)

		// wait for a signal
		<-signals

		// shut down the server
		ac.logger.Info("shutting down the server")
		shutdownError = h.Shutdown(context.Background())

		// log a shutdown error (factor into return value later)
		if shutdownError != nil {
			ac.logger.Error(shutdownError.Error())
		}
	} else if mode == "worker" {

		// wait for a signal
		<-signals

	} else {
		return errors.New("unknown nibbler mode detected: " + mode)
	}

	ac.logger.Info("shutting down the application")

	// destroy extensions in reverse order
	var err error
	var destroyError error
	for i := range ac.extensions {
		x := ac.extensions[len(ac.extensions)-i-1]
		destroyError = x.Destroy(ac)

		if destroyError != nil {
			err = destroyError
		}
	}

	// return a shutdown error if it occurred
	if shutdownError != nil {
		return shutdownError
	}

	// return any (the latest, which could be an improvement we could make) extension destroy error
	return err
}

func startServer(h *http.Server, ac *Application) error {

	// log that we're listening and state the port
	ac.logger.Info("Listening on " + strconv.Itoa((*ac.config).Port))

	// listen (this blocks) - log an error if it happened
	if err := h.ListenAndServe(); err != nil {
		ac.logger.Error("Failed to initialize server: " + err.Error())
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

type NoOpExtension struct {
}

func (s *NoOpExtension) Init(app *Application) error {
	return nil
}

func (s *NoOpExtension) Destroy(app *Application) error {
	return nil
}

func (s *NoOpExtension) AddRoutes(app *Application) error {
	return nil
}
