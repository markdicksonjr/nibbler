package session

import (
	"encoding/gob"
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/user"
)

type SessionStoreConnector interface {
	Connect() (error, sessions.Store)
}

type Extension struct {
	Secret         string
	SessionName    string
	StoreConnector *SessionStoreConnector	// creates cookie store if not provided
	store          *sessions.Store			// created by this extension
}

func (s *Extension) Init(app *nibbler.Application) error {
	if len(s.Secret) == 0 {
		return errors.New("session extension requires secret")
	}

	(*app.GetConfiguration().Raw).Get()

	gob.Register(map[string]interface{}{})

	// if a store connector is provided, use it
	if s.StoreConnector != nil {
		storeConnector := *s.StoreConnector
		errConnect, store := storeConnector.Connect()
		s.store = &store
		return errConnect
	}

	// if a connector isn't provided, use a cookie store
	var store sessions.Store = sessions.NewCookieStore([]byte(s.Secret))
	s.store = &store

	return nil
}

func (s *Extension) AddRoutes(app *nibbler.Application) error {
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return nil
}

func (s *Extension) GetAttribute(r *http.Request, attribute string) (interface{}, error) {
	session, err := (*s.store).Get(r, s.SessionName)

	if err != nil {
		return nil, err
	}

	sessionAttribute := session.Values[attribute]

	return sessionAttribute, nil
}

// TODO: SetAttributes to set multiple attributes in one save
func (s *Extension) SetAttribute(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	session, err := (*s.store).Get(r, s.SessionName)

	if err != nil {
		return err
	}

	session.Values[key] = value
	return session.Save(r, w)
}

func (s *Extension) GetCaller(r *http.Request) (*user.User, error) {
	sessionUser, err := s.GetAttribute(r, "user")

	if err != nil {
		return nil, err
	}

	if sessionUser == nil {
		return nil, nil
	}

	return user.FromJson(sessionUser.(string))
}

func (s *Extension) SetCaller(w http.ResponseWriter, r *http.Request, userValue *user.User) error {

	if userValue == nil {
		return s.SetAttribute(w, r, "user", nil)
	}

	userJson, err := user.ToJson(userValue)

	if err != nil {
		return err
	}

	return s.SetAttribute(w, r, "user", userJson)
}
