package repository

import (
	"go-api/dto/request/otpRequestDTO"
	"go-api/entity"

	"gorm.io/gorm"
)

type OTPRepository interface {
	UpdateOrCreate(data entity.TblOtp)
	FindByOTPAndStatus(request otpRequestDTO.ValidateOTPRequest, status string) entity.TblOtp
}

type otpConnection struct {
	connection *gorm.DB
}

func NewOTPRepository(conn *gorm.DB) OTPRepository {
	return &otpConnection{
		connection: conn,
	}
}

func (db *otpConnection) UpdateOrCreate(data entity.TblOtp) {
	var checker entity.TblOtp
	db.connection.Raw("SELECT * FROM tbl_otp where email = ?", data.Email).Take(&checker)
	if checker.Email == "" {
		db.connection.Save(&data)
	} else if checker.Email != "" {
		db.connection.Debug().Model(&entity.TblOtp{}).Where("email = ?", data.Email).Updates(&data)
	}
}

func (db *otpConnection) FindByOTPAndStatus(request otpRequestDTO.ValidateOTPRequest, status string) entity.TblOtp {
	var data entity.TblOtp
	db.connection.Model(&entity.TblOtp{}).Where("otp = ?", request.OTP).Where("status = ?", status).Take(&data)
	return data
}
