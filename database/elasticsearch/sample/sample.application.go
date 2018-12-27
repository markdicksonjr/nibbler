package main

import (
	"log"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/elasticsearch"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extension(s) for initialization
	extensions := []nibbler.Extension{
		&elasticsearch.Extension{},
	}

	// initialize the application
	appContext := nibbler.Application{}
	err = appContext.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = appContext.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
