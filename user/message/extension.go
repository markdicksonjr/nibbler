package message

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/session"
	"net/http"
	"strconv"
	"time"
)

type MessageType int

const (
	GENERAL MessageType = iota
	ALERT
	SYSTEM
)

type Message struct {
	ID           string      `json:"id" bson:"_id" gorm:"primary_key"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	DeletedAt    *time.Time  `json:"deletedAt,omitempty"`
	FromUserID   *string     `json:"fromUserId,omitempty"`
	FromUserName *string     `json:"fromUserName,omitempty"`
	Content      string      `json:"content"`
	Type         MessageType `json:"type"`
}

type UserMessageState struct {
	ID        string     `json:"id" bson:"_id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	UserID    string     `json:"userId"`
	ReadAt    *time.Time `json:"readAt,omitempty"`
}

type CompleteUserMessageState struct {
	UserMessageState
	Message
}

type PersistenceExtension interface {
	GetMessagesByUserId(userId string, count int, offset int) ([]CompleteUserMessageState, error)
	SendMessageToUser(userId string, message Message) error
	DeleteUserMessageState(messageStateId string, hardDelete bool) error
	MarkUserMessageStateAsRead(messageStateId string) error
}

type Extension struct {
	nibbler.NoOpExtension
	PersistenceExtension PersistenceExtension
	SessionExtension     *session.Extension
}

func (s *Extension) GetName() string {
	return "message"
}

func (s *Extension) PostInit(app *nibbler.Application) error {
	if s.PersistenceExtension == nil {
		return errors.New(s.GetName() + " requires a persistence extension but none was provided")
	}

	if s.SessionExtension == nil {
		return errors.New(s.GetName() + " requires a session extension but none was provided")
	}

	s.Logger = app.Logger

	app.Router.HandleFunc(app.Config.ApiPrefix + "/message", s.GetMessagesHandler).Queries("userId", "{userId}", "count", "{count}", "offset", "{offset}").Methods("GET")
	app.Router.HandleFunc(app.Config.ApiPrefix + "/message", s.SendMessageToUserHandler).Methods("POST")
	app.Router.HandleFunc(app.Config.ApiPrefix + "/message", s.DeleteUserMessageStateHandler).Queries("userId", "{userId}").Methods("DELETE")
	app.Router.HandleFunc(app.Config.ApiPrefix + "/message/read", s.MarkUserMessageStateAsReadHandler).Methods("POST")

	return nil
}

func (s *Extension) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	c, _ := strconv.Atoi(v["count"])
	o, _ := strconv.Atoi(v["offset"])

	// get the caller from the session so we can use it to authorize
	caller, err := s.SessionExtension.GetCaller(r)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}
	if caller == nil {
		nibbler.Write404Json(w)
		return
	}

	// currently, only a user can get their own messages
	if caller.ID != v["userId"] {
		s.Logger.Warn("a user with ID " + caller.ID + " requested mail for user with ID " + v["userId"])
		nibbler.Write404Json(w)
		return
	}

	// load the messages and states
	states, err := s.PersistenceExtension.GetMessagesByUserId(v["userId"], c, o)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	statesBytes, err := json.Marshal(states)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, string(statesBytes))
}

func (s *Extension) SendMessageToUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *Extension) DeleteUserMessageStateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *Extension) MarkUserMessageStateAsReadHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
