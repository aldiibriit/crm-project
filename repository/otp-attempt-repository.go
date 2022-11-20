package repository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type OTPAttemptRepository interface {
	UpdateOrCreate(data entity.TblOtpAttempt)
	FindByEmail(email string) entity.TblOtpAttempt
	Delete(email string)
}

type otpAttemptConnection struct {
	connection *gorm.DB
}

func NewOTPAttemptRepository(conn *gorm.DB) OTPAttemptRepository {
	return &otpAttemptConnection{
		connection: conn,
	}
}

func (db *otpAttemptConnection) UpdateOrCreate(data entity.TblOtpAttempt) {
	var checker entity.TblOtpAttempt
	db.connection.Debug().Raw("SELECT * FROM tbl_otp_attempt where email = ?", data.Email).Take(&checker)
	if checker.Email == "" {
		db.connection.Debug().Model(&entity.TblOtpAttempt{}).Create(&data)
	} else if checker.Email != "" {
		db.connection.Debug().Model(&entity.TblOtpAttempt{}).Where("email = ?", data.Email).Updates(&data)
	}
}

func (db *otpAttemptConnection) FindByEmail(email string) entity.TblOtpAttempt {
	var data entity.TblOtpAttempt
	db.connection.Model(&entity.TblOtpAttempt{}).Where("email = ?", email).Take(&data)
	return data
}

func (db *otpAttemptConnection) Delete(email string) {
	db.connection.Raw("DELETE FROM tbl_otp_attempt where email = ?", email)
}
