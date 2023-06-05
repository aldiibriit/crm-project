package entity

import "time"

// User represents users table in database
type User struct {
	ID       uint64 `gorm:"primary_key:auto_increment" json:"-"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Email    string `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password string `gorm:"->;<-;not null" json:"-"`
	Token    string `gorm:"-" json:"token,omitempty"`
}

type TblUser struct {
	ID             int    `gorm:"column:id" json:"-"`
	IdResponse     string `json:"id" gorm:"-"`
	Email          string
	MobileNo       string `json:"mobileNo" gorm:"column:mobile_no"`
	RegistrationId string
	Status         string
	Type           string
	Password       string    `json:"-"`
	CreatedAt      time.Time `json:"-"`
	ModifiedAt     time.Time `json:"-"`
	CreatedAtRes   string    `json:"createdAt" gorm:"-"`
	ModifiedAtRes  string    `json:"modifiedAt" gorm:"-"`
}
