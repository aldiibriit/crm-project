package repository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type EmailAttemptRepository interface {
	FindByEmailAndAction(email string, action int) entity.TblEmailAttempt
	Save(data entity.TblEmailAttempt)
	Update(data entity.TblEmailAttempt)
	UpdateOrCreate(data entity.TblEmailAttempt)
}

type emailAttempt struct {
	connection *gorm.DB
}

func NewEmailAttemptRepository(conn *gorm.DB) EmailAttemptRepository {
	return &emailAttempt{
		connection: conn,
	}
}

func (db *emailAttempt) FindByEmailAndAction(email string, action int) entity.TblEmailAttempt {
	var data entity.TblEmailAttempt
	db.connection.Debug().Model(&entity.TblEmailAttempt{}).Where("email = ? and action = ?", email, action).Find(&data)
	return data
}

func (db *emailAttempt) Save(data entity.TblEmailAttempt) {
	db.connection.Create(&data)
}

func (db *emailAttempt) Update(data entity.TblEmailAttempt) {
	db.connection.Model(&entity.TblEmailAttempt{}).Where("email = '?'", data.Email).Updates(&data)
}

func (db *emailAttempt) UpdateOrCreate(data entity.TblEmailAttempt) {
	var checker entity.TblEmailAttempt
	db.connection.Raw("SELECT * FROM tbl_email_attempt where email = ? and action = ?", data.Email, data.Action).Take(&checker)
	if checker.Email == "" {
		db.connection.Save(&data)
	} else if checker.Email != "" {
		db.connection.Debug().Model(&entity.TblEmailAttempt{}).Where("email = ?", data.Email).Where("action = ?", data.Action).Update("attempt", data.Attempt).Update("updated_at", data.UpdatedAt)
	}
}
