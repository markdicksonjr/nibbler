package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/sql"
	"log"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extensions for initialization
	sqlExtension := sql.Extension{
		//models = ...
	}
	extensions := []nibbler.Extension{
		&sqlExtension,
	}

	// initialize the application
	app := nibbler.Application{}
	if err = app.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
