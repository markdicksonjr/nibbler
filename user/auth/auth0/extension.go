package auth0

import (
	"net/http"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/auth/auth0"
	"github.com/markdicksonjr/nibbler/user"
)

type Extension struct {
	auth0.Extension

	UserExtension *user.Extension
}

func (s *Extension) Init(app *nibbler.Application) error {

	// call init on the base auth0 extension
	err := s.Extension.Init(app)

	// any error is fatal, return it
	if err != nil {
		return err
	}

	// apply a custom function for when we process the redirect from auth0
	s.OnLoginComplete = func(a *auth0.Extension, w http.ResponseWriter, r *http.Request) (allowRedirect bool, err error) {
		profile, err := s.SessionExtension.GetAttribute(r, "profile")

		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return false, err
		}

		profileMap := profile.(map[string]interface{})
		email := profileMap["name"]

		if email == nil {
			nibbler.Write500Json(w, err.Error())
			return false, err
		}

		emailValue := email.(string)

		if len(emailValue) == 0 {
			nibbler.Write500Json(w, err.Error())
			return false, err
		}

		userValue, err := s.UserExtension.GetUserByEmail(emailValue)

		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return false, err
		}

		if userValue == nil {
			nibbler.Write404Json(w)
			return false, err
		}

		err = s.SessionExtension.SetCaller(w, r, userValue)

		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return false, err
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "ok"}`))

		return false, nil
	}

	s.OnLogoutComplete = func(a *auth0.Extension, w http.ResponseWriter, r *http.Request) error {
		s.SessionExtension.SetCaller(w, r, nil)
		return nil
	}

	return nil
}

