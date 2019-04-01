package local

import (
	"github.com/google/uuid"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/mail/outbound"
	"github.com/markdicksonjr/nibbler/user"
	"net/http"
	"strings"
	"time"
)

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
		userValue, err = s.UserExtension.GetUserByEmail(email)
	} else if username != "" {
		userValue, err = s.UserExtension.GetUserByUsername(email)
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
		nibbler.Write200Json(w, `{"result": "ok"}`)
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

	go func() {

		// generate the link for the email
		var link = s.PasswordResetRedirect
		useAmpersand := strings.Contains(s.PasswordResetRedirect, "?")
		if useAmpersand {
			link += "&token=" + *userValue.PasswordResetToken
		} else {
			link += "?token=" + *userValue.PasswordResetToken
		}

		// build the recipient list
		emailVal := *userValue.Email
		var toList []*outbound.Email
		toList = append(toList, &outbound.Email{Address: emailVal, Name: name})

		// send the email
		_, err = s.Sender.SendMail(
			&outbound.Email{
				Name:    s.PasswordResetFromName,
				Address: s.PasswordResetFromEmail,
			},
			"Password Reset", // TODO: make configurable
			toList,
			"Please go to "+link+" to reset your password",
			"Please go to <a href=\""+link+"\">"+link+"</a> to reset your password",
		)
	}()

	nibbler.Write200Json(w, `{"result": "ok"}`)
}

func (s *Extension) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if !s.PasswordResetEnabled {
		nibbler.Write404Json(w)
		return
	}

	password := r.FormValue("password")
	token := r.FormValue("token")

	userValue, err := s.getUserByPasswordResetTokenAndValidate(token)

	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	if userValue == nil {
		nibbler.Write500Json(w, "please try again")
		return
	}

	// at this point, the token is verified

	// TODO: check password complexity...  make configurable...

	*userValue.Password = GeneratePasswordHash(password)
	(*userValue).PasswordResetToken = nil
	(*userValue).PasswordResetExpiration = nil

	// TODO: ensure extension sets above props to null, as well
	if err = s.UserExtension.UpdatePassword(userValue); err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, `{"result": "ok"}`)
}

func (s *Extension) getUserByPasswordResetTokenAndValidate(token string) (*user.User, error) {
	if !s.PasswordResetEnabled {
		return nil, nil
	}

	userValue, err := s.UserExtension.GetUserByPasswordResetToken(token)

	if err != nil || userValue == nil {
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
