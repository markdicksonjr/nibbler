package sql

import (
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/sql"
	"github.com/markdicksonjr/nibbler/user"
)

type Extension struct {
	nibbler.NoOpExtension
	SqlExtension *sql.Extension
}

func (s *Extension) Init(app *nibbler.Application) error {
	// sql extension AutoMigrates models
	return nil
}

func (s *Extension) GetUserById(id uint) (*user.User, error) {
	s.SqlExtension.Db.Error = nil

	userValue := user.User{}
	s.SqlExtension.Db = s.SqlExtension.Db.First(&userValue, id)

	if s.SqlExtension.Db.RecordNotFound() {
		return nil, nil
	}

	return &userValue, s.SqlExtension.Db.Error
}

func (s *Extension) GetUserByEmail(email string) (*user.User, error) {
	s.SqlExtension.Db.Error = nil

	userValue := user.User{}
	s.SqlExtension.Db = s.SqlExtension.Db.First(&userValue, "email = ?", email)

	if s.SqlExtension.Db.RecordNotFound() {
		return nil, nil
	}

	return &userValue, s.SqlExtension.Db.Error
}

func (s *Extension) GetUserByUsername(username string) (*user.User, error) {
	s.SqlExtension.Db.Error = nil

	userValue := user.User{}
	s.SqlExtension.Db = s.SqlExtension.Db.First(&userValue, "username = ?", username)

	if s.SqlExtension.Db.RecordNotFound() {
		return nil, nil
	}

	// TODO: nil, return code?, db error code?
	return &userValue, s.SqlExtension.Db.Error
}

func (s *Extension) GetUserByPasswordResetToken(token string) (*user.User, error) {
	s.SqlExtension.Db.Error = nil

	userValue := user.User{}
	s.SqlExtension.Db = s.SqlExtension.Db.First(&userValue, "password_reset_token = ?", token)

	if s.SqlExtension.Db.RecordNotFound() {
		return nil, nil
	}

	return &userValue, s.SqlExtension.Db.Error
}

func (s *Extension) GetUserByEmailValidationToken(token string) (*user.User, error) {
	s.SqlExtension.Db.Error = nil

	userValue := user.User{}
	s.SqlExtension.Db = s.SqlExtension.Db.First(&userValue, "email_validation_token = ?", token)

	if s.SqlExtension.Db.RecordNotFound() {
		return nil, nil
	}

	return &userValue, s.SqlExtension.Db.Error
}

func (s *Extension) Create(user *user.User) (*user.User, error) {
	s.SqlExtension.Db.Error = nil
	s.SqlExtension.Db = s.SqlExtension.Db.Create(user)
	// TODO: nil, return code?, db error code?
	return user, s.SqlExtension.Db.Error
}

func (s *Extension) Update(userValue *user.User) error {
	s.SqlExtension.Db.Error = nil
	s.SqlExtension.Db = s.SqlExtension.Db.Model(userValue).Updates(user.User{
		ID: userValue.ID,
		FirstName: userValue.FirstName,
		LastName: userValue.LastName,
		PasswordResetToken: userValue.PasswordResetToken,
		PasswordResetExpiration: userValue.PasswordResetExpiration,
	})
	return s.SqlExtension.Db.Error
}

func (s *Extension) UpdatePassword(userValue *user.User) (error) {
	s.SqlExtension.Db.Error = nil
	s.SqlExtension.Db = s.SqlExtension.Db.Model(userValue).Updates(user.User{
		ID: userValue.ID,
		Password: userValue.Password,
		PasswordResetToken: nil,
		PasswordResetExpiration: nil,
	})
	return s.SqlExtension.Db.Error
}