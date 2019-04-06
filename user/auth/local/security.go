package local

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// http://codahale.com/how-to-safely-store-a-password/

func GeneratePasswordHash(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)
	return string(hashedPwd), err
}

func ValidatePassword(password string, hashedPassword string) (bool, error) {
	if len(password) < 8 {
		return false, errors.New("password too short")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}
