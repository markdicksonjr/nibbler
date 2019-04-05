package main

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/socket/socket.io"
	"log"
)

func main() {

	// allocate configuration
	config, err := nibbler.LoadConfiguration(nil)
	if err != nil {
		log.Fatal(err)
	}

	socketIoExtension := socket_io.Extension{
		Port: 8000,
	}

	// initialize the application, provide config, logger, extensions
	appContext := nibbler.Application{}
	if err = appContext.Init(config, nibbler.DefaultLogger{}, []nibbler.Extension{
		&socketIoExtension,
	}); err != nil {
		log.Fatal(err.Error())
	}

	socketIoExtension.RegisterConnectHandler("test", "test", func(s socketio.Conn) error {
		log.Println("socket connect")
		return nil
	})

	// start the app
	if err = appContext.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
