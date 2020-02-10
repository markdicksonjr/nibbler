package local

import (
	"github.com/google/uuid"
	"github.com/markdicksonjr/nibbler"
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
	u, err := s.UserExtension.GetUserByEmail(email)
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
	userValue := nibbler.User{
		Email:            &email,
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
	encryptedPassword, err := GeneratePasswordHash(password)

	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

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
			var link = s.EmailVerificationRedirect
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
			var toList []*nibbler.EmailAddress
			toList = append(toList, &nibbler.EmailAddress{Address: emailVal, Name: name})

			// send the email
			_, err = s.Sender.SendMail(
				&nibbler.EmailAddress{
					Name:    s.EmailVerificationFromName,
					Address: s.EmailVerificationFromEmail,
				},
				"EmailAddress Verification", // TODO: make configurable
				toList,
				"Please go to "+link+" to verify your email",                          // TODO: configurable, with template param for link
				"Please go to <a href=\""+link+"\">"+link+"</a> to verify your email", // TODO: configurable, with template param for link
			)

			if err != nil {
				s.app.Logger.Error("while sending email verification, " + err.Error())
			}
		}()
	}

	if s.OnRegistrationSuccessful != nil {
		(*s.OnRegistrationSuccessful)(safeUser)
	}

	nibbler.Write200Json(w, `{"user": `+jsonString+`}`)
}

func (s *Extension) EmailTokenVerifyHandler(w http.ResponseWriter, r *http.Request) {

	// the endpoint is only available if verification is enabled
	if !s.EmailVerificationEnabled {
		nibbler.Write404Json(w)
		return
	}

	// grab and validate input parameters
	token := r.FormValue("token")
	if token == "" {
		nibbler.Write500Json(w, "a token form parameter is required")
		return
	}

	// look up the user from the token (factoring in expiration times)
	userValue, err := s.getUserByEmailValidationToken(token)

	// if an error happened during the lookup
	if err != nil {
		s.app.Logger.Error("while verifying email token: " + err.Error())
		nibbler.Write200Json(w, `{"result": false}`)
		return
	}

	// if no user has that email token
	if userValue == nil {
		nibbler.Write200Json(w, `{"result": false}`)
		return
	}

	// update the user in the DB
	isTrue := true
	userValue.IsEmailValidated = &isTrue
	userValue.EmailValidationToken = nil
	userValue.EmailValidationExpiration = nil
	if err = s.UserExtension.Update(userValue); err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	// grab the user from the session - it may not be there
	// token validation could happen while logged in, but will
	// more likely happen while not logged in
	sessionUser, err := s.SessionExtension.GetCaller(r)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	// if there's a user in the session, update the IsEmailValidated field of the user the session
	if sessionUser != nil {
		sessionUser.IsEmailValidated = &isTrue

		if err := s.SessionExtension.SetCaller(w, r, sessionUser); err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}
	}

	// if the success callback was provided, call it
	if s.OnEmailVerificationSuccessful != nil {
		safeUser := user.GetSafeUser(*userValue)
		(*s.OnEmailVerificationSuccessful)(safeUser)
	}

	nibbler.Write200Json(w, `{"result": true}`)
}

func (s *Extension) getUserByEmailValidationToken(token string) (*nibbler.User, error) {
	if !s.EmailVerificationEnabled {
		return nil, nil
	}

	userValue, err := s.UserExtension.GetUserByEmailVerificationToken(token)

	if err != nil || userValue == nil {
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
