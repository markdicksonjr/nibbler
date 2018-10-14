package main

import (
	"log"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/redis"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}
	config, err := nibbler.LoadApplicationConfiguration(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extensions for initialization
	redisExtension := redis.Extension{}
	extensions := []nibbler.Extension{
		&redisExtension,
	}

	// initialize the application
	app := nibbler.Application{}
	err = app.Init(config, &logger, &extensions)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}

