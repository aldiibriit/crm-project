package entity

import "time"

// User represents users table in database
type User struct {
	ID       uint64  `gorm:"primary_key:auto_increment" json:"-"`
	Name     string  `gorm:"type:varchar(255)" json:"name"`
	Email    string  `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password string  `gorm:"->;<-;not null" json:"-"`
	Token    string  `gorm:"-" json:"token,omitempty"`
	Books    *[]Book `json:"books,omitempty"`
}

type TblUser struct {
	ID             int `gorm:"column:id"`
	Email          string
	MobileNo       string
	RegistrationId string
	Status         string
	Type           string
	Password       string
	CreatedAt      time.Time
	ModifiedAt     time.Time
}
