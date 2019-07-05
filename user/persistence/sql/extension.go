package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/database/sql"
)

type Extension struct {
	nibbler.NoOpExtension
	SqlExtension *sql.Extension
}

func (s *Extension) Init(app *nibbler.Application) error {
	// sql extension AutoMigrates models
	return nil
}

func (s *Extension) GetUserById(id string) (*nibbler.User, error) {
	userValue := nibbler.User{}
	err := s.SqlExtension.Db.First(&userValue, id).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return &userValue, err
}

func (s *Extension) GetUserByEmail(email string) (*nibbler.User, error) {
	userValue := nibbler.User{}
	err := s.SqlExtension.Db.First(&userValue, "email = ?", email).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return &userValue, err
}

func (s *Extension) GetUserByUsername(username string) (*nibbler.User, error) {
	userValue := nibbler.User{}
	err := s.SqlExtension.Db.First(&userValue, "username = ?", username).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	// TODO: nil, return code?, db error code?
	return &userValue, err
}

func (s *Extension) GetUserByPasswordResetToken(token string) (*nibbler.User, error) {
	s.SqlExtension.Db.Error = nil

	userValue := nibbler.User{}
	err := s.SqlExtension.Db.First(&userValue, "password_reset_token = ?", token).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return &userValue, err
}

func (s *Extension) GetUserByEmailValidationToken(token string) (*nibbler.User, error) {
	s.SqlExtension.Db.Error = nil

	userValue := nibbler.User{}
	err := s.SqlExtension.Db.First(&userValue, "email_validation_token = ?", token).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return &userValue, err
}

func (s *Extension) Create(user *nibbler.User) (*nibbler.User, error) {
	err := s.SqlExtension.Db.Create(user).Error
	// TODO: nil, return code?, db error code?
	return user, err
}

func (s *Extension) Update(userValue *nibbler.User) error {
	// TODO: possibly use First(), update fields we care about, then use Save
	// Update will not save nil values, but Save will, presumably

	return s.SqlExtension.Db.Model(userValue).Updates(nibbler.User{
		ID:                      userValue.ID,
		FirstName:               userValue.FirstName,
		LastName:                userValue.LastName,
		PasswordResetToken:      userValue.PasswordResetToken,
		PasswordResetExpiration: userValue.PasswordResetExpiration,
	}).Error
}

func (s *Extension) UpdatePassword(userValue *nibbler.User) error {
	if err := s.SqlExtension.Db.Model(userValue).Updates(nibbler.User{
		ID:       userValue.ID,
		Password: userValue.Password,
	}).Error; err != nil {
		return err
	}

	if err := sql.NullifyField(s.SqlExtension.Db, "password_reset_token").Error; err != nil {
		return err
	}

	if err := sql.NullifyField(s.SqlExtension.Db, "password_reset_token_expiration").Error; err != nil {
		return err
	}

	return nil
}
