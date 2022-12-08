package repository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type KPRRepository interface {
	Insert(data entity.TblPengajuanKprBySales) entity.TblPengajuanKprBySales
}

type kprConnection struct {
	connection *gorm.DB
}

func NewKPRRepository(conn *gorm.DB) KPRRepository {
	return &kprConnection{
		connection: conn,
	}
}

func (db *kprConnection) Insert(data entity.TblPengajuanKprBySales) entity.TblPengajuanKprBySales {
	db.connection.Save(&data)
	return data
}
