package session

import (
	"encoding/gob"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/user"
	"net/http"
)

var requiresConnectorError = "session extension requires connector"

type StoreConnector interface {
	Connect() (error, sessions.Store) // TODO: reverse param order
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

	// if the extension hasn't had its logger set, take it from the app
	if s.Logger == nil {
		s.Logger = app.Logger
	}

	// if a store connector is provided, use it
	if s.StoreConnector == nil {
		return errors.New(requiresConnectorError)
	}

	errConnect, store := s.StoreConnector.Connect()

	// save the store to the extension
	s.store = &store
	return errConnect
}

func (s *Extension) GetName() string {
	return "session"
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

// SetCaller puts the user into the session as the current user
func (s *Extension) SetCaller(w http.ResponseWriter, r *http.Request, userValue *nibbler.User) error {

	if userValue == nil {
		return s.SetAttribute(w, r, "user", nil)
	}

	// wipe values for stringification TODO: use user.GetSafeString?
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

func (s *Extension) EnforceLoggedIn(routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.GetCaller(r)

		if err != nil {
			s.Logger.Error("while enforcing login, an error occurred: " + err.Error())
			nibbler.Write404Json(w)
			return
		}

		if caller == nil {
			nibbler.Write401Json(w)

			// TODO: allow optional callback for this
			s.Logger.Trace("received unauthorized request")
			return
		}

		routerFunc(w, r)
	}
}

// EnforceParamMatchesCallerID will validate that the user is logged in and is requesting a resource that's identified by
// a user ID that matches the caller's ID.  The param argument to this function determines which path param is checked
// against the caller's ID.
func (s *Extension) EnforceParamMatchesCallerID(param string, routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.GetCaller(r)
		if err != nil {
			s.Logger.Error("while params matches caller ID, an error occurred: " + err.Error())
			nibbler.Write404Json(w)
			return
		}

		if caller == nil {
			nibbler.Write401Json(w)
			s.Logger.Trace("received unauthorized request")
			return
		}

		if mux.Vars(r)[param] != caller.ID {
			nibbler.Write404Json(w)
			s.Logger.Trace("received unauthorized request from caller with ID " + caller.ID + " against param " + param + " with value " + mux.Vars(r)[param])
			return
		}

		routerFunc(w, r)
	}
}

// EnforceEmailValidated validates the user is logged in and that user's email has been validated.  Note that the session
// extension can't validate emails, but extensions like the user-local-auth extension can
func (s *Extension) EnforceEmailValidated(routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.GetCaller(r)
		if err != nil {
			s.Logger.Error("while enforcing email validated, an error occurred: " + err.Error())
			nibbler.Write404Json(w)
			return
		}

		if caller == nil {
			nibbler.Write404Json(w)
			// TODO: log
			return
		}

		if caller.IsEmailValidated == nil || !*caller.IsEmailValidated {
			nibbler.Write404Json(w)
			return
		}

		routerFunc(w, r)
	}
}