package session

import (
	"encoding/gob"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/user"
	"net/http"
)

type StoreConnector interface {
	Connect() (error, sessions.Store)
	MaxAge() int
}

type Extension struct {
	nibbler.NoOpExtension
	SessionName    string
	StoreConnector StoreConnector  // creates cookie store if not provided
	store          *sessions.Store // created by this extension
}

func (s *Extension) Init(app *nibbler.Application) error {
	gob.Register(map[string]interface{}{})

	// if a store connector is provided, use it
	if s.StoreConnector == nil {
		return errors.New("session extension requires connector")
	}

	errConnect, store := s.StoreConnector.Connect()

	// save the store to the extension
	s.store = &store
	return errConnect
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

func (s *Extension) GetCaller(r *http.Request) (*nibbler.User, error) {
	sessionUser, err := s.GetAttribute(r, "user")

	if err != nil {
		return nil, err
	}

	if sessionUser == nil {
		return nil, nil
	}

	return user.FromJson(sessionUser.(string))
}

func (s *Extension) SetCaller(w http.ResponseWriter, r *http.Request, userValue *nibbler.User) error {

	if userValue == nil {
		return s.SetAttribute(w, r, "user", nil)
	}

	// wipe values for stringification
	password := userValue.Password
	resetToken := userValue.PasswordResetToken
	resetExpiration := userValue.PasswordResetExpiration
	userValue.Password = nil
	userValue.PasswordResetToken = nil
	userValue.PasswordResetExpiration = nil

	userJson, err := user.ToJson(userValue)

	// restore password and reset token
	userValue.Password = password
	userValue.PasswordResetToken = resetToken
	userValue.PasswordResetExpiration = resetExpiration

	if err != nil {
		return err
	}

	return s.SetAttribute(w, r, "user", userJson)
}
