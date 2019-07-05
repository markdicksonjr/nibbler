package local

import (
	"errors"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/user"
	"net/http"
)

type Extension struct {
	nibbler.Extension

	SessionExtension *session.Extension
	UserExtension    *user.Extension

	// for emailing
	Sender nibbler.MailSender

	// for password reset
	PasswordResetEnabled             bool
	PasswordResetFromName            string
	PasswordResetFromEmail           string
	PasswordResetRedirect            string // a UI or other service to handle the redirect from email (will have ?token=X or &token=X appended)
	PasswordResetTokenExpirationDays *int

	// for email verification
	RegistrationEnabled                  bool
	EmailVerificationEnabled             bool // whether email verification is available (doesn't mean it's required)
	EmailVerificationRequired            bool // whether email verification is required before logging in
	EmailVerificationTokenExpirationDays *int
	EmailVerificationRedirect            string
	EmailVerificationFromName            string
	EmailVerificationFromEmail           string

	// callbacks (for extending default behavior)
	OnLoginSuccessful             *func(loggedInUser nibbler.User, sessionMaxAgeMinutes int)
	OnLogoutSuccessful            *func(loggedOutUser nibbler.User)
	OnRegistrationSuccessful      *func(registeredUser nibbler.User)
	OnEmailVerificationSuccessful *func(registeredUser nibbler.User)

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

	// if registration is enabled, check prerequisites
	if s.EmailVerificationEnabled {
		if s.Sender == nil {
			return errors.New("sender extension was not provided to user local auth extension, but features using it are enabled")
		}

		if s.EmailVerificationFromName == "" {
			return errors.New("email verification from name was not provided to user local auth extension, but features using it are enabled")
		}

		if s.EmailVerificationFromEmail == "" {
			return errors.New("email verification from address was not provided to user local auth extension, but features using it are enabled")
		}
	}

	return nil
}

func (s *Extension) AddRoutes(app *nibbler.Application) error {
	app.GetRouter().HandleFunc("/api/user", s.GetCurrentUserHandler).Methods("GET")
	app.GetRouter().HandleFunc("/api/login", s.LoginFormHandler).Methods("POST")
	app.GetRouter().HandleFunc("/api/logout", s.LogoutHandler).Methods("POST", "GET")
	app.GetRouter().HandleFunc("/api/password/reset-token", s.ResetPasswordTokenHandler).Methods("POST")
	app.GetRouter().HandleFunc("/api/password", s.ResetPasswordHandler).Methods("POST")

	if s.RegistrationEnabled {
		app.GetRouter().HandleFunc("/api/register", s.RegisterFormHandler).Methods("POST")

		if s.EmailVerificationEnabled {
			app.GetRouter().HandleFunc("/api/email/validate", s.EmailTokenVerifyHandler).Methods("POST")
		}
	}
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return nil
}

func (s *Extension) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, err := s.SessionExtension.GetCaller(r)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	if currentUser == nil {
		nibbler.Write404Json(w)
		return
	}

	safeUser := user.GetSafeUser(*currentUser)
	jsonString, err := user.ToJson(&safeUser)

	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, jsonString)
}
