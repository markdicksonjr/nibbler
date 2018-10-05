package user

import (
	"encoding/json"
	"time"
)

type User struct {
	ID						string		`json:"id" bson:"_id" gorm:"primary_key"`	// ID
	CreatedAt				time.Time	`json:"createdAt"`
	UpdatedAt				time.Time	`json:"updatedAt"`
	DeletedAt				*time.Time	`json:"deletedAt,omitempty" sql:"index"`
	Email					*string		`json:"email,omitempty" gorm:"unique"`
	Username				*string		`json:"username,omitempty" gorm:"unique"`
	FirstName				*string		`json:"firstName,omitempty"`
	LastName 				*string		`json:"lastName,omitempty"`
	Title					*string		`json:"title,omitempty"`
	Password  				*string		`json:"password,omitempty"`
	PortraitUri 			*string		`json:"portraitUri,omitempty"`
	AvatarUri				*string		`json:"avatarUri,omitempty"`
	StatusText  			*string		`json:"statusText,omitempty"`
	IsActive				bool		`json:"isActive"`
	IsEmailValidated		bool		`json:"isEmailValidated"`
	DeactivatedAt			*time.Time  `json:"deactivatedAt,omitempty"`
	LastLogin				*time.Time	`json:"lastLogin,omitempty"`
	FailedLoginCount		int8		`json:"failedLoginCount"`
	Gender					*string		`json:"gender,omitempty" gorm:"size:1"`
	PhoneHome				*string		`json:"phoneHome,omitempty" gorm:"size:24"`
	PhoneWork				*string		`json:"phoneWork,omitempty" gorm:"size:24"`
	PhoneMobile				*string		`json:"phoneMobile,omitempty" gorm:"size:24"`
	PhoneOther				*string		`json:"phoneOther,omitempty" gorm:"size:24"`
	Fax						*string		`json:"fax,omitempty" gorm:"size:24"`
	Uri						*string		`json:"uri,omitempty"`
	Birthday				*string		`json:"birthday,omitempty"`
	Bio						*string		`json:"bio,omitempty"`
	AddressLine1			*string		`json:"addressLine1,omitempty"`
	AddressLine2			*string		`json:"addressLine2,omitempty"`
	AddressLine3			*string		`json:"addressLine3,omitempty"`
	AddressCity				*string		`json:"addressCity,omitempty"`
	AddressStateOrProvince	*string		`json:"addressStateOrProvince,omitempty"`
	AddressPostalCode		*string		`gorm:"size:16" json:"postalCode,omitempty"`
	CountryCode				*string		`gorm:"size:3" json:"countryCode,omitempty"`
	EmployeeId				*string		`json:"employeeId,omitempty"`
	ReferenceId				*string		`json:"referenceId,omitempty"`
	PasswordResetToken		*string		`json:"passwordResetToken,omitempty"`
	PasswordResetExpiration	*time.Time	`json:"passwordResetExpiration,omitempty"`
	EmploymentStartDate		*time.Time	`json:"employmentStartDate,omitempty"`
	EmploymentEndDate		*time.Time	`json:"employmentEndDate,omitempty"`
	ContractStartDate		*time.Time	`json:"contractStartDate,omitempty"`
	ContractEndDate			*time.Time	`json:"contractEndDate,omitempty"`
	Context					*string		`json:"contractEndDate,omitempty"`
}

func FromJson(jsonString string) (*User, error) {
	byteArray := []byte(jsonString)
	user := User{}
	err := json.Unmarshal(byteArray, &user)
	return &user, err
}

func ToJson(user *User) (result string, err error) {
	userJsonBytes, err := json.Marshal(user)

	if err != nil {
		return
	}

	result = string(userJsonBytes)

	return
}