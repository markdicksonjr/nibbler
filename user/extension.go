package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/markdicksonjr/nibbler"
)

const noExtensionErrorMessage = "no extension found"

type PersistenceController interface {
	Init(app *nibbler.Application) error
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
	
	PersistenceController PersistenceController
}

func (s *Extension) Init(app *nibbler.Application) error {
	if s.PersistenceController == nil {
		return errors.New("persistence controller was not provided to user extension")
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
	if s.PersistenceController != nil {
		return s.PersistenceController.GetUserById(id)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByEmail(email string) (*User, error) {
	if s.PersistenceController != nil {
		return s.PersistenceController.GetUserByEmail(email)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByPasswordResetToken(token string) (*User, error) {
	if s.PersistenceController != nil {
		return s.PersistenceController.GetUserByPasswordResetToken(token)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByUsername(username string) (*User, error) {
	if s.PersistenceController != nil {
		return s.PersistenceController.GetUserByUsername(username)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) Create(user *User) (*User, error) {
	if s.PersistenceController != nil {
		user.ID = uuid.New().String()
		return s.PersistenceController.Create(user)
	}
	return user, errors.New(noExtensionErrorMessage)
}

func (s *Extension) Update(user *User) (error) {
	if s.PersistenceController != nil {
		return s.PersistenceController.Update(user)
	}
	return errors.New(noExtensionErrorMessage)
}

func (s *Extension) UpdatePassword(user *User) (error) {
	if s.PersistenceController != nil {
		return s.PersistenceController.UpdatePassword(user)
	}
	return errors.New(noExtensionErrorMessage)
}
