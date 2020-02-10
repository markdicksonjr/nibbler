package user

import (
	"encoding/json"
	"github.com/markdicksonjr/nibbler"
)

func FromJson(jsonString string) (*nibbler.User, error) {
	u := nibbler.User{}
	if err := json.Unmarshal([]byte(jsonString), &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func ToJson(user *nibbler.User) (result string, err error) {
	var userJsonBytes []byte
	userJsonBytes, err = json.Marshal(user)
	if err != nil {
		return
	}

	result = string(userJsonBytes)
	return
}

func GetSafeUser(user nibbler.User) nibbler.User {
	safeUser := user
	safeUser.Password = nil
	safeUser.PasswordResetExpiration = nil
	safeUser.PasswordResetToken = nil
	safeUser.EmailValidationToken = nil
	safeUser.EmailValidationExpiration = nil
	safeUser.ProtectedContext = nil
	return safeUser
}
