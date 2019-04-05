package main

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/redis"
	"log"
	"time"
)

func main() {

	// allocate logger and configuration
	config, err := nibbler.LoadConfiguration(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// prepare extensions for initialization
	redisExtension := redis.Extension{}

	// initialize the application, provide config, logger, extensions
	app := nibbler.Application{}
	if err = app.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&redisExtension,
	}); err != nil {
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

	if err = app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
