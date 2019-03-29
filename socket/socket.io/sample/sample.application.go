package main

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/socket/socket.io"
	"log"
)

func main() {

	// allocate logger and configuration
	var logger nibbler.Logger = nibbler.DefaultLogger{}

	// allocate configuration
	config, err := nibbler.LoadConfiguration(nil)
	if err != nil {
		log.Fatal(err)
	}

	socketIoExtension := socket_io.Extension{
		Port: 8000,
	}

	// prepare extensions for initialization
	extensions := []nibbler.Extension{
		&socketIoExtension,
	}

	// initialize the application
	appContext := nibbler.Application{}
	if err = appContext.Init(config, &logger, &extensions); err != nil {
		log.Fatal(err.Error())
	}

	socketIoExtension.RegisterConnectHandler("test", "test", func(s socketio.Conn) error {
		logger.Debug("socket connect")
		return nil
	})

	// start the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
