package entity

import "time"

type TblOtp struct {
	ID         string    `json:"id"`
	OTP        string    `json:"otp"`
	MobileNo   string    `json:"mobileNo"`
	Email      string    `json:"email"`
	Status     string    `json:"status"`
	OtpAttempt int       `json:"otpAttempt"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	ExpiredAt  time.Time `json:"expiredAt"`
}
