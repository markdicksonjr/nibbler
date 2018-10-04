package nibbler

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/micro/go-config"
)

type Extension interface {
	Init(app *Application) error
	AddRoutes(app *Application) error
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

	// set up the static directory routing
	http.Handle(configValue.StaticDirectory + "/", http.StripPrefix(configValue.StaticDirectory, http.FileServer(http.Dir(configValue.StaticDirectory))))

	// allocate a router for everything else
	ac.router = mux.NewRouter()

	// init extension routes
	for _, x := range extensionValue {
		err = x.AddRoutes(ac)

		if err != nil {
			return err
		}
	}

	loggerValue.Info("Starting server")

	http.Handle("/", ac.router)

	return nil
}

func (ac *Application) Run() error {

	// dereference parameters for ease-of-use
	configValue := *ac.config
	loggerValue := *ac.logger

	// log that we're listening and state the port
	loggerValue.Info("Listening on " + strconv.Itoa(configValue.Port))

	// listen (this blocks)
	err := http.ListenAndServe(":" + strconv.Itoa(configValue.Port), nil)

	// log an error if it happened
	if err != nil {
		loggerValue.Error("Failed to initialize server: " + err.Error());
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