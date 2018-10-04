package user

import (
	"encoding/json"
	"time"
)

type User struct {
	ID						string		`json:"id" bson:"_id" gorm:"primary_key"`	// ID
	CreatedAt				time.Time	`json:"createdAt"`
	UpdatedAt				time.Time	`json:"updatedAt"`
	DeletedAt				*time.Time	`json:"deletedAt" sql:"index"`
	Email					*string		`json:"email" gorm:"unique"`
	Username				*string		`json:"username" gorm:"unique"`
	FirstName				*string		`json:"firstName"`
	LastName 				*string		`json:"lastName"`
	Title					*string		`json:"title"`
	Password  				*string		`json:"password"`
	PortraitUri 			*string		`json:"portraitUri"`
	AvatarUri				*string		`json:"avatarUri"`
	StatusText  			*string		`json:"statusText"`
	IsActive				bool		`json:"isActive"`
	IsEmailValidated		bool		`json:"isEmailValidated"`
	DeactivatedAt			*time.Time  `json:"deactivatedAt"`
	LastLogin				*time.Time	`json:"lastLogin"`
	FailedLoginCount		int8		`json:"failedLoginCount"`
	Gender					*string		`json:"gender" gorm:"size:1"`
	PhoneHome				*string		`json:"phoneHome" gorm:"size:24"`
	PhoneWork				*string		`json:"phoneWork" gorm:"size:24"`
	PhoneMobile				*string		`json:"phoneMobile" gorm:"size:24"`
	PhoneOther				*string		`json:"phoneOther" gorm:"size:24"`
	Fax						*string		`json:"fax" gorm:"size:24"`
	Uri						*string		`json:"uri"`
	Birthday				*string		`json:"birthday"`
	Bio						*string		`json:"bio"`
	AddressLine1			*string		`json:"addressLine1"`
	AddressLine2			*string		`json:"addressLine2"`
	AddressLine3			*string		`json:"addressLine3"`
	AddressCity				*string		`json:"addressCity"`
	AddressStateOrProvince	*string		`json:"addressStateOrProvince"`
	AddressPostalCode		*string		`gorm:"size:16" json:"postalCode"`
	CountryCode				*string		`gorm:"size:3" json:"countryCode"`
	EmployeeId				*string		`json:"employeeId"`
	ReferenceId				*string		`json:"referenceId"`
	PasswordResetToken		*string		`json:"passwordResetToken"`
	PasswordResetExpiration	*time.Time	`json:"passwordResetExpiration"`
	EmploymentStartDate		*time.Time	`json:"employmentStartDate"`
	EmploymentEndDate		*time.Time	`json:"employmentEndDate"`
	ContractStartDate		*time.Time	`json:"contractStartDate"`
	ContractEndDate			*time.Time	`json:"contractEndDate"`
	Context					*string		`json:"contractEndDate"`
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