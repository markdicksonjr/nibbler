package local

import (
	"errors"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/user"
	"net/http"
	"strconv"
)

func (s *Extension) EnforceLoggedIn(routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.SessionExtension.GetCaller(r)

		if err != nil {
			s.app.Logger.Error("while enforcing login, an error occurred: " + err.Error())
			nibbler.Write404Json(w)
			return
		}

		if caller == nil {
			nibbler.Write404Json(w)
			// TODO: log
			return
		}

		routerFunc(w, r)
	}
}

// also validates the user is logged in
func (s *Extension) EnforceEmailValidated(routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.SessionExtension.GetCaller(r)

		if err != nil {
			s.app.Logger.Error("while enforcing email validated, an error occurred: " + err.Error())
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

func (s *Extension) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	userValue, err := s.Login(email, password)

	// if an error happened during login
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	// if the user isn't in the system
	if userValue == nil {
		nibbler.Write500Json(w, "please try again")
		return
	}

	// set the caller in the session
	if err = s.SessionExtension.SetCaller(w, r, userValue); err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	safeUser := user.GetSafeUser(*userValue)
	jsonString, err := user.ToJson(&safeUser)

	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	if s.OnLoginSuccessful != nil {
		(*s.OnLoginSuccessful)(safeUser, s.SessionExtension.StoreConnector.MaxAge())
	}

	nibbler.Write200Json(w, `{"user": `+jsonString+
		`, "sessionAgeSeconds":`+strconv.Itoa(s.SessionExtension.StoreConnector.MaxAge())+`}`)
}

func (s *Extension) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := s.SessionExtension.SetCaller(w, r, nil)

	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	if s.OnLogoutSuccessful != nil {
		sessionUser, err := s.SessionExtension.GetCaller(r)

		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		(*s.OnLogoutSuccessful)(*sessionUser)
	}

	nibbler.Write200Json(w, `{"result": "ok"}`)
}

func (s *Extension) Login(email string, password string) (*nibbler.User, error) {
	u, err := s.UserExtension.GetUserByEmail(email)
	if err != nil {
		s.app.Logger.Error("while looking up user by email, error = " + err.Error())
		return u, err
	}

	if u == nil || u.Password == nil {
		return nil, nil
	}

	validPassword, err := ValidatePassword(password, *u.Password)
	if err != nil {
		s.app.Logger.Error("while validating password in login flow, error = " + err.Error())
		return nil, err
	}

	if !validPassword {
		s.app.Logger.Trace("invalid password for email " + email)
		return nil, errors.New("invalid password")
	}

	// if we need email verification but it hasn't been done yet, fail
	if s.EmailVerificationEnabled && s.EmailVerificationRequired && (u.IsEmailValidated == nil || !*u.IsEmailValidated) {
		s.app.Logger.Debug("login blocked for email " + email + " because it was not verified")
		return nil, errors.New("email not verified")
	}

	return u, nil
}
