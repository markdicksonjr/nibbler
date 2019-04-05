package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/elasticsearch"
	"log"
)

func main() {

	// allocate logger and configuration
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize the application, provide config, logger, extensions
	appContext := nibbler.Application{}
	if err = appContext.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&elasticsearch.Extension{},
	}); err != nil {
		log.Fatal(err.Error())
	}

	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
