package repository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	Insert(data entity.TblCustomer) entity.TblCustomer
}

type customerConnection struct {
	connection *gorm.DB
}

func NewCustomerRepository(conn *gorm.DB) CustomerRepository {
	return &customerConnection{
		connection: conn,
	}
}

func (db *customerConnection) Insert(data entity.TblCustomer) entity.TblCustomer {
	db.connection.Save(&data)
	return data
}
