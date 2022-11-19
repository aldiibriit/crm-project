package repository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type SalesRepository interface {
	InsertRelation(data entity.TblSales) error
}

type salesConnection struct {
	connection *gorm.DB
}

func NewSalesRepository(conn *gorm.DB) SalesRepository {
	return &salesConnection{
		connection: conn,
	}
}

func (db *salesConnection) InsertRelation(data entity.TblSales) error {
	err := db.connection.Save(&data).Error
	if err != nil {
		return err
	}
	return nil
}
