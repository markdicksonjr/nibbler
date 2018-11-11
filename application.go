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
	Init(app *Application)		error
	AddRoutes(app *Application)	error
	Destroy(app *Application)	error
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
	Raw                 *config.Config
	StaticDirectory     string
}

type HeaderConfiguration struct {
	AccessControlAllowHeaders	string
	AccessControlAllowMethods	string
	AccessControlAllowOrigin	string
}

type Application struct {
	config		*Configuration
	extensions	*[]Extension
	logger		*Logger
	router		*mux.Router
}

func (ac *Application) Init(config *Configuration, logger *Logger, extensions *[]Extension) error {
	ac.config = config
	ac.logger = logger
	ac.extensions = extensions

	// dereference parameters for ease-of-use
	extensionValue := *extensions
	configValue := *config
	loggerValue := *logger

	// prepare a general-use error variable
	var err error = nil

	// initialize all extensions
	for _, x := range extensionValue {
		err = x.Init(ac)

		// if any error occurred, return the error and stop processing
		if err != nil {
			return err
		}
	}

	// allocate a router
	ac.router = mux.NewRouter()

	// init extension routes
	for _, x := range extensionValue {
		err = x.AddRoutes(ac)

		if err != nil {
			return err
		}
	}

	// set up the static directory routing
	ac.router.PathPrefix("/").Handler(http.FileServer(http.Dir(configValue.StaticDirectory)))


	loggerValue.Info("Starting server")

	http.Handle("/", ac.router)

	return nil
}

func (ac *Application) Run() error {

	// dereference parameters for ease-of-use
	loggerValue := *ac.logger

	// get the configured app mode, accounting for the default value
	mode := (*ac.config.Raw).Get("nibbler", "mode").String("web")

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
		loggerValue.Info("shutting down the server")
		shutdownError = h.Shutdown(context.Background())

		// log a shutdown error (factor into return value later)
		if shutdownError != nil {
			loggerValue.Error(shutdownError.Error())
		}
	} else if mode == "worker" {

		// wait for a signal
		<-signals

	} else {
		return errors.New("unknown nibbler mode detected: " + mode)
	}

	loggerValue.Info("shutting down the application")

	// dereference parameters for ease-of-use
	extensionValue := *ac.extensions

	// destroy extensions in reverse order
	var err error
	var destroyError error
	for i := range extensionValue {
		x := extensionValue[len(extensionValue) - i - 1]
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

	// dereference parameters for ease-of-use
	loggerValue := *ac.logger

	// log that we're listening and state the port
	loggerValue.Info("Listening on " + strconv.Itoa((*ac.config).Port))

	// listen (this blocks)
	err := h.ListenAndServe()

	// log an error if it happened
	if err != nil {
		loggerValue.Error("Failed to initialize server: " + err.Error())
	}

	return err
}

func (ac *Application) GetLogger() *Logger {
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