package nibbler

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
	"os"
)

func GetConfigurationFromSources(sources []source.Source) (config.Config, error) {
	conf := config.NewConfig()

	// load each source
	for _, sourceItem := range sources {
		err := conf.Load(sourceItem)

		if err != nil {
			return nil, err
		}
	}

	return conf, nil
}

// "merging priority is in reverse order"
// if nil, environment source used
func LoadConfiguration(sources *[]source.Source) (*Configuration, error) {

	// if sources are not provided
	if sources == nil {
		var envSources []source.Source

		configFileExists := true
		if _, err := os.Stat("./config.json"); err != nil {
			if os.IsNotExist(err) {
				configFileExists = false
			}
		}

		// allow a file to override env config, if it exists
		if configFileExists {
			envSources = append(envSources, file.NewSource(file.WithPath("./config.json")))
		}

		// use environment variable source
		envSources = append(envSources, env.NewSource())

		sources = &envSources
	}

	// load the app configuration
	conf, err := GetConfigurationFromSources(*sources)

	// if an error occurred, return it
	if err != nil {
		return nil, err
	}

	// get NIBBLER_PORT and PORT, giving precendence to NIBBLER_PORT
	// PORT is a common PaaS requirement to even have the app run
	primaryPort := conf.Get("nibbler", "port").Int(3000)
	secondaryPort := conf.Get("port").Int(3000)

	if primaryPort == 3000 && secondaryPort != primaryPort {
		primaryPort = secondaryPort
	}

	return &Configuration{
		Raw:             &conf,
		Port:            primaryPort,
		StaticDirectory: conf.Get("nibbler", "directory", "static").String("./public/"),
		HeaderConfiguration: HeaderConfiguration{
			AccessControlAllowOrigin: conf.Get("nibbler", "ac", "allow", "origin").String("*"),
			AccessControlAllowMethods: conf.Get("nibbler", "ac", "allow", "methods").String("GET, POST, OPTIONS, PUT, PATCH, DELETE"),
			AccessControlAllowHeaders: conf.Get("nibbler", "ac", "allow", "headers").String("Origin, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, X-Response-Time, X-PINGOTHER, X-CSRF-Token, Authorization"),
		},
	}, nil
}
