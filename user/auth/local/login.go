package local

import (
	"encoding/json"
	"errors"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/user"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (s *Extension) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if email == "" && username == "" && password == "" && r.Body != nil {
		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if bodyBytes == nil {
			nibbler.Write500Json(w, "{\"message\": \"body was not json\"")
			return
		}

		var asMap map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &asMap); err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		email, _ = asMap["email"].(string)
		email = strings.TrimSpace(email)

		username, _ = asMap["username"].(string)
		username = strings.TrimSpace(username)

		password, _ = asMap["password"].(string)
		password = strings.TrimSpace(password)
	}

	userValue, err := s.Login(email, username, password)

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

func (s *Extension) Login(email string, username string, password string) (*nibbler.User, error) {
	var u *nibbler.User
	var err error

	if email != "" {
		u, err = s.UserExtension.GetUserByEmail(email)
		if err != nil {
			s.app.Logger.Error("while looking up user by email, error = " + err.Error())
			return u, err
		}
	} else if username != "" {
		u, err = s.UserExtension.GetUserByUsername(username)
		if err != nil {
			s.app.Logger.Error("while looking up user by usernae, error = " + err.Error())
			return u, err
		}
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
		s.app.Logger.Trace("invalid password for email \"" + email + "\", username \"" + username + "\"")
		return nil, errors.New("invalid password")
	}

	// if we need email verification but it hasn't been done yet, fail
	if s.EmailVerificationEnabled && s.EmailVerificationRequired && (u.IsEmailValidated == nil || !*u.IsEmailValidated) {
		s.app.Logger.Debug("login blocked for email " + email + " because it was not verified")
		return nil, errors.New("email not verified")
	}

	return u, nil
}
