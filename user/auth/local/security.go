package local

import (
	"errors"
	"github.com/jameskeane/bcrypt"
)

const HASH_LENGTH = 29

// http://codahale.com/how-to-safely-store-a-password/

func GeneratePasswordHash(password string) string {
	// generate a random salt with default rounds of complexity (10)
	salt, _ := bcrypt.Salt()

	// hash and verify a password with random salt
	hash, _ := bcrypt.Hash(password, salt)

	return hash + salt
}

func ValidatePassword(password string, hashedPassword string) (bool, error) {
	if len(password) < 8 {
		return false, errors.New("password too short")
	}
	if len(hashedPassword) < HASH_LENGTH {
		return false, errors.New("stored password is invalid")
	}

	hashSalt := hashedPassword[len(hashedPassword)-HASH_LENGTH : len(hashedPassword)]

	hash, _ := bcrypt.Hash(password, hashSalt)
	hashedTest := hash + hashSalt

	if hashedTest == hashedPassword {
		return true, nil
	}
	return false, errors.New("password did not match")
}