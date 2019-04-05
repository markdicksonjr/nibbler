package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/sql"
	"log"
)

func main() {

	// allocate logger and configuration
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extensions for initialization
	sqlExtension := sql.Extension{
		//models = ...
	}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&sqlExtension,
	}); err != nil {
		log.Fatal(err.Error())
	}

	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
