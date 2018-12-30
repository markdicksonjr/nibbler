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
			// TODO: log
			w.WriteHeader(http.StatusNotFound)
			//w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`404 page not found`))
			return
		}

		if caller == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`404 page not found`))
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

	nibbler.Write200Json(w, `{"user": ` + jsonString +
		`, "sessionAgeSeconds":`+ strconv.Itoa(s.SessionExtension.StoreConnector.MaxAge()) + `}`)
}

func (s *Extension) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	s.SessionExtension.SetCaller(w, r, nil)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "ok"}`))
}

func (s *Extension) Login(email string, password string) (*user.User, error) {
	u, err := s.UserExtension.GetUserByEmail(email);

	if err != nil {
		return u, err
	}

	if u == nil || u.Password == nil {
		return nil, nil
	}

	validPassword, err := ValidatePassword(password, *u.Password);

	if err != nil {
		return nil, err
	}

	if !validPassword {
		return nil, errors.New("invalid password")
	}

	return u, nil;
}