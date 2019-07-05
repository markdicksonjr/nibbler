package user

import (
	"encoding/json"
	"github.com/markdicksonjr/nibbler"
	"reflect"
)

func FromJson(jsonString string) (*nibbler.User, error) {
	userInt, err := nibbler.FromJson(jsonString, reflect.TypeOf(nibbler.User{}))
	return userInt.(*nibbler.User), err
}

func ToJson(user *nibbler.User) (result string, err error) {
	userJsonBytes, err := json.Marshal(user)

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
