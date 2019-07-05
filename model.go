package nibbler

import (
	"time"
)

type User struct {
	ID                        string     `json:"id" bson:"_id" gorm:"primary_key"` // ID
	CreatedAt                 time.Time  `json:"createdAt"`
	UpdatedAt                 time.Time  `json:"updatedAt"`
	DeletedAt                 *time.Time `json:"deletedAt,omitempty" sql:"index"`
	Email                     *string    `json:"email,omitempty" gorm:"unique"`
	Username                  *string    `json:"username,omitempty" gorm:"unique"`
	FirstName                 *string    `json:"firstName,omitempty"`
	LastName                  *string    `json:"lastName,omitempty"`
	Title                     *string    `json:"title,omitempty"`
	Password                  *string    `json:"password,omitempty"`
	PortraitUri               *string    `json:"portraitUri,omitempty"`
	AvatarUri                 *string    `json:"avatarUri,omitempty"`
	StatusText                *string    `json:"statusText,omitempty"`
	IsActive                  *bool      `json:"isActive,omitempty"`
	IsEmailValidated          *bool      `json:"isEmailValidated,omitempty"`
	DeactivatedAt             *time.Time `json:"deactivatedAt,omitempty"`
	LastLogin                 *time.Time `json:"lastLogin,omitempty"`
	FailedLoginCount          *int8      `json:"failedLoginCount,omitempty"`
	Gender                    *string    `json:"gender,omitempty" gorm:"size:1"`
	PhoneHome                 *string    `json:"phoneHome,omitempty" gorm:"size:24"`
	PhoneWork                 *string    `json:"phoneWork,omitempty" gorm:"size:24"`
	PhoneMobile               *string    `json:"phoneMobile,omitempty" gorm:"size:24"`
	PhoneOther                *string    `json:"phoneOther,omitempty" gorm:"size:24"`
	Fax                       *string    `json:"fax,omitempty" gorm:"size:24"`
	Uri                       *string    `json:"uri,omitempty"`
	Birthday                  *string    `json:"birthday,omitempty"`
	Bio                       *string    `json:"bio,omitempty"`
	AddressLine1              *string    `json:"addressLine1,omitempty"`
	AddressLine2              *string    `json:"addressLine2,omitempty"`
	AddressLine3              *string    `json:"addressLine3,omitempty"`
	AddressCity               *string    `json:"addressCity,omitempty"`
	AddressStateOrProvince    *string    `json:"addressStateOrProvince,omitempty"`
	AddressPostalCode         *string    `gorm:"size:16" json:"postalCode,omitempty"`
	CountryCode               *string    `gorm:"size:3" json:"countryCode,omitempty"`
	EmployeeId                *string    `json:"employeeId,omitempty"`
	ReferenceId               *string    `json:"referenceId,omitempty"`
	PasswordResetToken        *string    `json:"passwordResetToken,omitempty"`
	PasswordResetExpiration   *time.Time `json:"passwordResetExpiration,omitempty"`
	EmailValidationToken      *string    `json:"emailValidationToken,omitempty"`
	EmailValidationExpiration *time.Time `json:"emailValidationExpiration,omitempty"`
	EmploymentStartDate       *time.Time `json:"employmentStartDate,omitempty"`
	EmploymentEndDate         *time.Time `json:"employmentEndDate,omitempty"`
	ContractStartDate         *time.Time `json:"contractStartDate,omitempty"`
	ContractEndDate           *time.Time `json:"contractEndDate,omitempty"`
	PrimaryLocation           *string    `json:"primaryLocation,omitempty"` // e.g. lat/long, grid codes, etc
	Context                   *string    `json:"context,omitempty"`
	ProtectedContext          *string    `json:"protectedContext,omitempty"`
}

// basic model for both role-based and group privilege-based auth control

type Group struct {
	ID        string     `json:"id" bson:"_id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" sql:"index"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
}

type GroupMembership struct {
	ID        string     `json:"id" bson:"_id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" sql:"index"`
	GroupID   string     `json:"groupId"`
	MemberID  string     `json:"memberId"`
	Role      GroupRole  `json:"role"`
}

type GroupPrivilege struct {
	ID                string     `json:"id" bson:"_id" gorm:"primary_key"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt,omitempty" sql:"index"`
	PerformingGroupID string     `json:"performingGroupID"` // e.g. "administrators" ID
	TargetGroupID     string     `json:"targetGroupID"`     // e.g. "customers" ID
	Action            int        `json:"action"`            // e.g. read/write/admin/etc
}

type GroupRole int // make your own roles
type Action int    // make your own actions