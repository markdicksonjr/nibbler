package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/markdicksonjr/nibbler"
)

const noExtensionErrorMessage = "no extension found"

type PersistenceExtension interface {
	nibbler.Extension
	GetUserById(id string) (*nibbler.User, error)
	GetUserByEmail(email string) (*nibbler.User, error)
	GetUserByUsername(username string) (*nibbler.User, error)
	Create(user *nibbler.User) (*nibbler.User, error)
	Update(user *nibbler.User) error
	UpdatePassword(user *nibbler.User) error
	GetUserByPasswordResetToken(token string) (*nibbler.User, error)
	GetUserByEmailValidationToken(token string) (*nibbler.User, error)
}

type Extension struct {
	nibbler.Extension

	PersistenceExtension PersistenceExtension

	OnBeforeUserCreate     func(user *nibbler.User)
	OnAfterUserCreate      func(user *nibbler.User)
	OnBeforeUserUpdate     func(user *nibbler.User)
	OnAfterUserUpdate      func(user *nibbler.User)
	OnBeforePasswordUpdate func(user *nibbler.User)
	OnAfterPasswordUpdate  func(user *nibbler.User)
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

func (s *Extension) GetUserById(id string) (*nibbler.User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserById(id)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByEmail(email string) (*nibbler.User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserByEmail(email)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByPasswordResetToken(token string) (*nibbler.User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserByPasswordResetToken(token)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByEmailVerificationToken(token string) (*nibbler.User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserByEmailValidationToken(token)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) GetUserByUsername(username string) (*nibbler.User, error) {
	if s.PersistenceExtension != nil {
		return s.PersistenceExtension.GetUserByUsername(username)
	}
	return nil, errors.New(noExtensionErrorMessage)
}

func (s *Extension) Create(user *nibbler.User) (*nibbler.User, error) {
	if s.PersistenceExtension != nil {
		user.ID = uuid.New().String()

		// call the OnBeforeUserCreate callback if provided
		if s.OnBeforeUserCreate != nil {
			s.OnBeforeUserCreate(user)
		}

		// save the new user to the DB
		resultUser, err := s.PersistenceExtension.Create(user)

		// if an error occurred, return now
		if err != nil {
			return resultUser, err
		}

		// call the OnAfterUserCreate callback if provided
		if s.OnAfterUserCreate != nil {
			s.OnAfterUserCreate(resultUser)
		}

		return resultUser, err
	}
	return user, errors.New(noExtensionErrorMessage)
}

func (s *Extension) Update(user *nibbler.User) error {
	if s.PersistenceExtension != nil {

		// call the OnBeforeUserUpdate callback if provided
		if s.OnBeforeUserUpdate != nil {
			s.OnBeforeUserUpdate(user)
		}

		// change user user in the DB
		if err := s.PersistenceExtension.Update(user); err != nil {
			return err
		}

		// call the OnAfterUserUpdate callback if provided
		if s.OnAfterUserUpdate != nil {
			s.OnAfterUserUpdate(user)
		}

		return nil
	}
	return errors.New(noExtensionErrorMessage)
}

func (s *Extension) UpdatePassword(user *nibbler.User) error {
	if s.PersistenceExtension != nil {

		// call the OnBeforePasswordUpdate callback if provided
		if s.OnBeforePasswordUpdate != nil {
			s.OnBeforePasswordUpdate(user)
		}

		// change user user in the DB
		if err := s.PersistenceExtension.UpdatePassword(user); err != nil {
			return err
		}

		// call the OnAfterUserUpdate callback if provided
		if s.OnAfterPasswordUpdate != nil {
			s.OnAfterPasswordUpdate(user)
		}

		return nil
	}
	return errors.New(noExtensionErrorMessage)
}
