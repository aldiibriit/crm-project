package InternalRepository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type LogActivityRepository interface {
	Insert(data entity.TbLogActivity) error
	InsertWithTx(data entity.TbLogActivity, tx *gorm.DB) error
}

type logActivityConnection struct {
	connection *gorm.DB
}

func NewLogActivityRepository(conn *gorm.DB) LogActivityRepository {
	return &logActivityConnection{
		connection: conn,
	}
}

func (db *logActivityConnection) Insert(data entity.TbLogActivity) error {
	return db.connection.Debug().Create(&data).Error
}

func (db *logActivityConnection) InsertWithTx(data entity.TbLogActivity, tx *gorm.DB) error {
	return tx.Debug().Create(&data).Error
}
