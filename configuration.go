package nibbler

import (
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/env"
	"github.com/micro/go-micro/config/source/file"
	"os"
)

// GetConfigurationFromSources will handle loading a configuration from a provided list of sources
func GetConfigurationFromSources(sources []source.Source) (config.Config, error) {
	conf := config.NewConfig()

	// load each source
	for _, sourceItem := range sources {
		if err := conf.Load(sourceItem); err != nil {
			return nil, err
		}
	}

	return conf, nil
}

// LoadConfiguration will handle loading the configuration from multiple sources, allowing for this module's default
// source set to be used if none are provided to this function
//
// "merging priority is in reverse order"
// if nil or empty, file and environment sources used (file takes precedence)
func LoadConfiguration(sources ...source.Source) (*Configuration, error) {

	// if sources are not provided
	if len(sources) == 0 {
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
		sources = envSources
	}

	// load the app configuration
	conf, err := GetConfigurationFromSources(sources)

	// if an error occurred, return it
	if err != nil {
		return nil, err
	}

	// get NIBBLER_PORT and PORT, giving precendence to NIBBLER_PORT
	// PORT is a common PaaS requirement to even have the app run
	primaryPort := conf.Get("nibbler", "port").Int(0)
	secondaryPort := conf.Get("port").Int(0)

	if primaryPort == 0 && secondaryPort != primaryPort {
		primaryPort = secondaryPort
	}

	return &Configuration{
		Raw:             conf,
		Port:            primaryPort,
		StaticDirectory: conf.Get("nibbler", "directory", "static").String("./public/"),
		ApiPrefix:       conf.Get("nibbler", "api", "prefix").String("/api"),
		Headers: HeaderConfiguration{
			AccessControlAllowOrigin:  conf.Get("nibbler", "ac", "allow", "origin").String("*"),
			AccessControlAllowMethods: conf.Get("nibbler", "ac", "allow", "methods").String("GET, POST, OPTIONS, PUT, PATCH, DELETE"),
			AccessControlAllowHeaders: conf.Get("nibbler", "ac", "allow", "headers").String("Origin, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, X-MailSendResponse-Time, X-PINGOTHER, X-CSRF-Token, Authorization"),
		},
	}, nil
}
