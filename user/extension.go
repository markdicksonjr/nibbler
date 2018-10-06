package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/markdicksonjr/nibbler"
)

const noExtensionErrorMessage = "no extension found"

type PersistenceExtension interface {
	nibbler.Extension
	GetUserById(id uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	Create(user *User) (*User, error)
	Update(user *User) (error)
	UpdatePassword(user *User) (error)
	GetUserByPasswordResetToken(token string) (*User, error)
}

type Extension struct {
	nibbler.Extension

	PersistenceExtension PersistenceExtension
}

func (s *Extension) Init(app *nibbler.Application) error {
	if s.PersistenceExtension == nil {
		return errors.New("no persistence extension was provided to user extension")
	}

	return nil
}

func (s *Extension) AddRoutes(app *nibbler.Application) error {
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return nil
}

func (s *Extension) GetSafeUser(user User) User {
	safeUser := user
	safeUser.Password = nil
	safeUser.PasswordResetExpiration = nil
	safeUser.PasswordResetToken = nil
	return safeUser
}

func (s *Extension) GetUserById(id uint) (*User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserById(id)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByEmail(email string) (*User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserByEmail(email)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByPasswordResetToken(token string) (*User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserByPasswordResetToken(token)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByUsername(username string) (*User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserByUsername(username)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) Create(user *User) (*User, error) {
	if s.PersistenceExtension != nil {
		user.ID = uuid.New().String()
		return s.PersistenceExtension.Create(user)
	}
	return user, errors.New(noExtensionErrorMessage)
}

func (s *Extension) Update(user *User) (error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.Update(user)
	}
	return errors.New(noExtensionErrorMessage)
}

func (s *Extension) UpdatePassword(user *User) (error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.UpdatePassword(user)
	}
	return errors.New(noExtensionErrorMessage)
}
