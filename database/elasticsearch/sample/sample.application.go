package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/elasticsearch"
	"log"
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
	if err = appContext.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
