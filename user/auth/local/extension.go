package local

import (
	"errors"
	"net/http"
	"github.com/markdicksonjr/nibbler"
	NibUser "github.com/markdicksonjr/nibbler/user"
	NibSendGrid "github.com/markdicksonjr/nibbler/mail/sendgrid"
	NibSession "github.com/markdicksonjr/nibbler/session"
	"github.com/google/uuid"
	"time"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"strings"
	_ "github.com/lib/pq" // TODO: does this go here?  doubtful
)

type Extension struct {
	SessionExtension *NibSession.Extension
	UserExtension *NibUser.Extension

	// for emailing
	SendGridExtension *NibSendGrid.Extension

	// for passwo
	//[[constraint]]
	//  revision = "bb90f23ce7de089adeda8ef6586a25b38ffbfaa2"
	//  name = "github.com/markdicksonjr/nibble"rd reset
	PasswordResetEnabled   bool
	PasswordResetFromName  string
	PasswordResetFromEmail string
	PasswordResetRedirect  string // a UI or other service to handle the redirect from email (will have ?token=X or &token=X appended)
}

func (s *Extension) Init(app *nibbler.Application) error {

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
		if s.SendGridExtension == nil {
			return errors.New("sendgrid extension was not provided to user local auth extension, but features using it are enabled")
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

	user, err := s.Login(email, password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "` + err.Error() + `"}`))
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "please try again"}`))
		return
	}

	s.SessionExtension.SetCaller(w, r, user)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "ok"}`))
}

func (s *Extension) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	s.SessionExtension.SetCaller(w, r, nil)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "ok"}`))
}

func (s *Extension) Login(email string, password string) (*NibUser.User, error) {
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
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "not found"}`)) // TODO: ensure this conforms
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")

	var user *NibUser.User
	var err error

	if email != "" {
		user, err = s.UserExtension.GetUserByEmail(email);
	} else if username != "" {
		user, err = s.UserExtension.GetUserByUsername(email);
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "incorrect parameters"}`))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "` + err.Error() + `"}`))
		return
	}

	// user not found, but we want to respond with OK to not give away too much info (i.e. our user list could be brute-forced)
	if user == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "ok"}`))
		return
	}

	// in the event we looked up the user by anything but email, check that there is an email
	if user.Email == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "no email on record"}`))
		return
	}

	// generate reset token with expiration (one day - TODO: make configurable)
	uuidInstance := uuid.New().String()
	expiration := time.Now().AddDate(0, 0, 1)
	user.PasswordResetToken = &uuidInstance
	user.PasswordResetExpiration = &expiration

	errUpdate := s.UserExtension.Update(user)

	if errUpdate != nil {
		// TODO: log error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "failed to update user record"}`))
		return
	}

	name := ""
	if user.FirstName != nil && user.LastName != nil {
		name = *user.FirstName + " " + *user.LastName
	}

	// generate the link for the email
	var link = s.PasswordResetRedirect
	useAmpersand := strings.Contains(s.PasswordResetRedirect, "?")
	if useAmpersand {
		link += "&token=" + *user.PasswordResetToken
	} else {
		link += "?token=" + *user.PasswordResetToken
	}

	go func() {

		// send email
		emailVal := *user.Email
		_, err = s.SendGridExtension.SendMail(
			mail.NewEmail(s.PasswordResetFromName, s.PasswordResetFromEmail),
			"Password Reset", // TODO: make configurable
			mail.NewEmail(name, emailVal),
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
	user, err := s.getUserByValidToken(token)

	if err != nil {
		// TODO: log err?
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": false}`))
		return
	}

	if user == nil {
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

	user, err := s.getUserByValidToken(token)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "` + err.Error() + `"}`))
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "please try again"}`))
		return
	}

	// at this point, the token is verified

	// TODO: check password complexity...

	*user.Password = GeneratePasswordHash(password)
	err = s.UserExtension.UpdatePassword(user)

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

func (s *Extension) getUserByValidToken(token string) (*NibUser.User, error) {
	if !s.PasswordResetEnabled {
		return nil, nil
	}

	user, err := s.UserExtension.GetUserByPasswordResetToken(token)

	if err != nil || user == nil{
		return nil, nil
	}

	if user.PasswordResetExpiration == nil {
		return nil, nil
	}

	if !(*user.PasswordResetExpiration).After(time.Now()) {
		return nil, nil
	}

	return user, nil
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
