package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/redis"
	"log"
	"time"
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

	cmd := redisExtension.Client.Set("test", "sd", time.Minute)
	if cmd.Err() != nil {
		log.Fatal(cmd.Err())
	}

	strCmd := redisExtension.Client.Get("test")
	if strCmd.Err() != nil {
		log.Fatal(strCmd.Err())
	}

	log.Println(strCmd.Val() == "sd")

	err = app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}

