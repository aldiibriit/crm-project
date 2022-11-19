package entity

import "time"

type TblEmailAttempt struct {
	Id         string    `json:"id" gorm:"size:255;column:id"`
	Action     int       `json:"action"`
	Attempt    int       `json:"attempt"`
	Email      string    `json:"email" gorm:"size:255"`
	UrlEncoded string    `json:"urlEncoded" gorm:"type:longtext"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
