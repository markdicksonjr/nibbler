package auth0

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/auth/auth0"
	"github.com/markdicksonjr/nibbler/user"
	"net/http"
)

type Extension struct {
	auth0.Extension

	UserExtension *user.Extension
}

func (s *Extension) Init(app *nibbler.Application) error {

	// call init on the base auth0 extension
	if err := s.Extension.Init(app); err != nil {
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
			nibbler.Write500Json(w, "no email was found for this profile")
			return false, err
		}

		emailValue := email.(string)

		if len(emailValue) == 0 {
			nibbler.Write500Json(w, "a blank email was found for this profile")
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

		if err = s.SessionExtension.SetCaller(w, r, userValue); err != nil {
			nibbler.Write500Json(w, err.Error())
			return false, err
		}

		if len(s.LoggedInRedirectUrl) > 0 {
			http.Redirect(w, r, s.LoggedInRedirectUrl, http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"result": "ok"}`))

		return false, err
	}

	s.OnLogoutComplete = func(a *auth0.Extension, w http.ResponseWriter, r *http.Request) error {
		return s.SessionExtension.SetCaller(w, r, nil)
	}

	return nil
}
