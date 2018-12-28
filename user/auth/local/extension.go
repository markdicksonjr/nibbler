package local

import (
	"errors"
	"github.com/google/uuid"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/user"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Extension struct {
	nibbler.Extension

	SessionExtension *session.Extension
	UserExtension *user.Extension

	// for emailing
	Sender outbound.Sender

	// for password reset
	PasswordResetEnabled   				bool
	PasswordResetFromName  				string
	PasswordResetFromEmail 				string
	PasswordResetRedirect  				string // a UI or other service to handle the redirect from email (will have ?token=X or &token=X appended)
	PasswordResetTokenExpirationDays	*int

	app *nibbler.Application
}

func (s *Extension) Init(app *nibbler.Application) error {
	s.app = app

	// assert that we have the session extension
	if s.SessionExtension == nil {
		return errors.New("session extension was not provided to user local auth extension")
	}

	// assert that we have the user extension
	if s.UserExtension == nil {
		return errors.New("user extension was not provided to user local auth extension")
	}

	// if password reset is enabled, check prerequisites
	if s.PasswordResetEnabled {
		if s.Sender == nil {
			return errors.New("sender extension was not provided to user local auth extension, but features using it are enabled")
		}

		if s.PasswordResetFromName == "" {
			return errors.New("password reset from name was not provided to user local auth extension, but features using it are enabled")
		}

		if s.PasswordResetFromEmail == "" {
			return errors.New("password reset from address was not provided to user local auth extension, but features using it are enabled")
		}
	}

	return nil
}

func (s *Extension) AddRoutes(app *nibbler.Application) error {
	app.GetRouter().HandleFunc("/api/login", s.LoginFormHandler).Methods("POST")
	app.GetRouter().HandleFunc("/api/logout", s.LogoutHandler).Methods("POST", "GET")
	app.GetRouter().HandleFunc("/api/password/reset-token", s.ResetPasswordTokenHandler).Methods("POST")
	app.GetRouter().HandleFunc("/api/password/reset-token/validity", s.ResetPasswordTokenVerifyHandler).Methods("POST")
	app.GetRouter().HandleFunc("/api/password", s.ResetPasswordHandler).Methods("POST")
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return nil
}

func (s *Extension) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	userValue, err := s.Login(email, password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "` + err.Error() + `"}`))
		return
	}

	if userValue == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "please try again"}`))
		return
	}

	s.SessionExtension.SetCaller(w, r, userValue)

	safeUser := user.GetSafeUser(*userValue)
	jsonString, err := user.ToJson(&safeUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"user": ` + jsonString +
		`, "sessionAgeMs"":`+ strconv.Itoa(s.SessionExtension.StoreConnector.MaxAge()) + `}`))
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

func (s *Extension) ResetPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	if !s.PasswordResetEnabled {
		nibbler.Write404Json(w)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")

	var userValue *user.User
	var err error

	if email != "" {
		userValue, err = s.UserExtension.GetUserByEmail(email);
	} else if username != "" {
		userValue, err = s.UserExtension.GetUserByUsername(email);
	} else {
		nibbler.Write500Json(w, "incorrect parameters")
		return
	}

	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	// user not found, but we want to respond with OK to not give away too much info (i.e. our user list could be brute-forced)
	if userValue == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "ok"}`))
		return
	}

	// in the event we looked up the user by anything but email, check that there is an email
	if userValue.Email == nil {
		nibbler.Write500Json(w, "no email on record")
		return
	}

	// compute password expiration time (defaults to 1 day)
	expirationDays := 1
	if s.PasswordResetTokenExpirationDays != nil {
		expirationDays = *s.PasswordResetTokenExpirationDays
	}

	// generate reset token with expiration
	uuidInstance := uuid.New().String()
	expiration := time.Now().AddDate(0, 0, expirationDays)
	userValue.PasswordResetToken = &uuidInstance
	userValue.PasswordResetExpiration = &expiration

	errUpdate := s.UserExtension.Update(userValue)

	if errUpdate != nil {
		(*s.app.GetLogger()).Error("failed to update user record: " + errUpdate.Error())
		nibbler.Write500Json(w, "failed to update user record")
		return
	}

	name := ""
	if userValue.FirstName != nil && userValue.LastName != nil {
		name = *userValue.FirstName + " " + *userValue.LastName
	}

	// generate the link for the email
	var link = s.PasswordResetRedirect
	useAmpersand := strings.Contains(s.PasswordResetRedirect, "?")
	if useAmpersand {
		link += "&token=" + *userValue.PasswordResetToken
	} else {
		link += "?token=" + *userValue.PasswordResetToken
	}

	go func() {

		// build the recipient list
		emailVal := *userValue.Email
		var toList []*outbound.Email
		toList = append(toList, &outbound.Email{Address: emailVal, Name: name})

		// send the email
		_, err = s.Sender.SendMail(
			&outbound.Email{
				Name: s.PasswordResetFromName,
				Address: s.PasswordResetFromEmail,
			},
			"Password Reset", // TODO: make configurable
			toList,
			"Please go to " + link + " to reset your password",
			"Please go to <a href=\"" + link + "\">" + link + "</a> to reset your password",
		)
	}()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "ok"}`))
}

func (s *Extension) ResetPasswordTokenVerifyHandler(w http.ResponseWriter, r *http.Request) {
	if !s.PasswordResetEnabled {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "not found"}`)) // TODO: ensure this conforms
		return
	}

	token := r.FormValue("token")
	userValue, err := s.getUserByValidToken(token)

	if err != nil {
		// TODO: log err?
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": false}`))
		return
	}

	if userValue == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": false}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": true}`))
}

func (s *Extension) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	token := r.FormValue("token")

	userValue, err := s.getUserByValidToken(token)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "` + err.Error() + `"}`))
		return
	}

	if userValue == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "please try again"}`))
		return
	}

	// at this point, the token is verified

	// TODO: check password complexity...

	*userValue.Password = GeneratePasswordHash(password)
	err = s.UserExtension.UpdatePassword(userValue)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "ok"}`))
}

func (s *Extension) getUserByValidToken(token string) (*user.User, error) {
	if !s.PasswordResetEnabled {
		return nil, nil
	}

	userValue, err := s.UserExtension.GetUserByPasswordResetToken(token)

	if err != nil || userValue == nil{
		return nil, nil
	}

	if userValue.PasswordResetExpiration == nil {
		return nil, nil
	}

	if !(*userValue.PasswordResetExpiration).After(time.Now()) {
		return nil, nil
	}

	return userValue, nil
}

// TODO: endpoint for setting the new password - re-confirm the token before allowing

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
