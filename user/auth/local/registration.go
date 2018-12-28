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

// TODO: allow username
func (s *Extension) RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	// enforce that a password is provided
	if password == "" {
		nibbler.Write500Json(w, "password is a required field")
		return
	}

	// TODO: password meets requirements

	// look up the user with that email
	u, err := s.UserExtension.GetUserByEmail(email);
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	// if the user is found
	// TODO: improve error message, but don't let the user know the email is in the system, if possible
	if u != nil {
		nibbler.Write500Json(w, "please try again")
		return
	}

	// begin putting together a new user
	emailValidated := !s.EmailVerificationEnabled
	userValue := user.User{
		Email: &email,
		IsEmailValidated: &emailValidated,
	}

	if s.EmailVerificationEnabled {

		// compute verification token expiration time (defaults to 1 day)
		expirationDays := 1
		if s.EmailVerificationTokenExpirationDays != nil {
			expirationDays = *s.EmailVerificationTokenExpirationDays
		}

		// generate verification token with expiration
		uuidInstance := uuid.New().String()
		expiration := time.Now().AddDate(0, 0, expirationDays)
		userValue.EmailValidationToken = &uuidInstance
		userValue.EmailValidationExpiration = &expiration
	}

	// compute and set the encrypted password
	encryptedPassword := GeneratePasswordHash(password)
	userValue.Password = &encryptedPassword

	// create a test user, if it does not exist
	u, err = s.UserExtension.Create(&userValue)

	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	// prepare to return the newly-created user to the client
	// we want to ensure the sensitive fields aren't in the response
	safeUser := user.GetSafeUser(*u)
	jsonString, err := user.ToJson(&safeUser)

	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	if s.EmailVerificationEnabled {

		// send email to verify the email for the account
		go func() {

			// generate the link for the email
			var link= s.EmailVerificationRedirect
			useAmpersand := strings.Contains(s.EmailVerificationRedirect, "?")
			if useAmpersand {
				link += "&token=" + *userValue.EmailValidationToken
			} else {
				link += "?token=" + *userValue.EmailValidationToken
			}

			name := ""
			if userValue.FirstName != nil && userValue.LastName != nil {
				name = *userValue.FirstName + " " + *userValue.LastName
			}

			// build the recipient list
			emailVal := *userValue.Email
			var toList []*outbound.Email
			toList = append(toList, &outbound.Email{Address: emailVal, Name: name})

			// send the email
			_, err = s.Sender.SendMail(
				&outbound.Email{
					Name: s.EmailVerificationFromName,
					Address: s.EmailVerificationFromEmail,
				},
				"Email Verification", // TODO: make configurable
				toList,
				"Please go to " + link + " to verify your email", // TODO: configurable, with template param for link
				"Please go to <a href=\"" + link + "\">" + link + "</a> to verify your email", // TODO: configurable, with template param for link
			)

			if err != nil {
				(*s.app.GetLogger()).Error("while sending email verification, " + err.Error())
			}
		}()
	}

	nibbler.Write200Json(w, `{"user": ` + jsonString + `}`)
}

func (s *Extension) EmailTokenVerifyHandler(w http.ResponseWriter, r *http.Request) {
	if !s.EmailVerificationEnabled {
		nibbler.Write404Json(w)
		return
	}

	token := r.FormValue("token")
	if token == "" {
		nibbler.Write500Json(w, "a token form parameter is required")
		return
	}

	userValue, err := s.getUserByEmailValidationToken(token)

	if err != nil {
		// TODO: log err?
		nibbler.Write200Json(w, `{"result": false}`)
		return
	}

	if userValue == nil {
		nibbler.Write200Json(w, `{"result": false}`)
		return
	}

	isTrue := true
	userValue.IsEmailValidated = &isTrue
	userValue.EmailValidationToken = nil
	userValue.EmailValidationExpiration = nil
	if err = s.UserExtension.Update(userValue); err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, `{"result": true}`)
}

func (s *Extension) getUserByEmailValidationToken(token string) (*user.User, error) {
	if !s.EmailVerificationEnabled {
		return nil, nil
	}

	userValue, err := s.UserExtension.GetUserByEmailVerificationToken(token)

	if err != nil || userValue == nil{
		return nil, nil
	}

	if userValue.EmailValidationToken == nil {
		return nil, nil
	}

	if !(*userValue.EmailValidationExpiration).After(time.Now()) {
		return nil, nil
	}

	return userValue, nil
}