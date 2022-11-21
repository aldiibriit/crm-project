package repository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type SalesRepository interface {
	InsertRelation(data entity.TblSales) error
	FindByEmailDeveloper(emailDeveloper string) []entity.TblSales
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

func (db *salesConnection) FindByEmailDeveloper(emailDeveloper string) []entity.TblSales {
	var result []entity.TblSales
	db.connection.Raw("SELECT * FROM tbl_sales WHERE developer_email = ?", emailDeveloper).Find(&result)
	return result
}
