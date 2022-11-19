package entity

import "time"

type TblOtpAttempt struct {
	Id         string    `json:"id"`
	Email      string    `json:"email"`
	OtpAttempt int       `json:"otpAttempt"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
