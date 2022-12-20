package repository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type KPRRepository interface {
	Insert(data entity.TblPengajuanKprBySales) entity.TblPengajuanKprBySales
	Delete(id string) entity.TblPengajuanKprBySales
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

func (db *kprConnection) Delete(id string) entity.TblPengajuanKprBySales {
	var data entity.TblPengajuanKprBySales
	db.connection.Model(&entity.TblPengajuanKprBySales{}).Where("id", id).Update("status", "on_deleted").Find(&data)
	return data
}
