package socket_io

import (
	"github.com/googollee/go-socket.io"
	"github.com/markdicksonjr/nibbler"
	"log"
	"net/http"
	"strconv"
)

type Extension struct {
	nibbler.NoOpExtension
	server *socketio.Server

	Port int
}

func (s *Extension) Init(app *nibbler.Application) error {
	var err error
	s.server, err = socketio.NewServer(nil)
	if err != nil {
		return err
	}

	go s.server.Serve()
	go func() {
		http.Handle("/socket.io/", s.server)
		if err := http.ListenAndServe(":" + strconv.Itoa(s.Port), nil); err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return s.server.Close()
}

type fnEventMessageHandler func(s socketio.Conn, msg string) string
type fnEventMessageHandlerVoid func(s socketio.Conn, msg string)
type fnEventWithErrorHandler func(s socketio.Conn) error
type fnErrorHandler func(err error)

func (s *Extension) RegisterEventHandler(ns, event string, handler fnEventMessageHandler) {
	s.server.OnEvent(ns, event, handler)
}

func (s *Extension) RegisterConnectHandler(ns, event string, handler fnEventWithErrorHandler) {
	s.server.OnConnect(ns, handler)
}

func (s *Extension) RegisterDisconnectHandler(ns, event string, handler fnEventMessageHandlerVoid) {
	s.server.OnDisconnect(ns, handler)
}

func (s *Extension) RegisterErrorHandler(ns string, handler fnErrorHandler) {
	s.server.OnError(ns, handler)
}
