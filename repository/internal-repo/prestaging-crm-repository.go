package InternalRepository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type PrestagingCRMRepository interface {
	Insert(data entity.TbPrestagingCrm) error
	InsertWithTx(data entity.TbPrestagingCrm, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
}

type prestagingCRMConnection struct {
	connection *gorm.DB
}

func NewPrestagingCRMRepository(conn *gorm.DB) PrestagingCRMRepository {
	return &prestagingCRMConnection{
		connection: conn,
	}
}

func (db *prestagingCRMConnection) Insert(data entity.TbPrestagingCrm) error {
	return db.connection.Debug().Create(&data).Error
}

func (db *prestagingCRMConnection) InsertWithTx(data entity.TbPrestagingCrm, tx *gorm.DB) error {
	return tx.Debug().Create(&data).Error
}

func (db *prestagingCRMConnection) UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error {
	return tx.Debug().Model(&entity.TbPrestagingCrm{}).Where("sn", data["sn"].(string)).Updates(&data).Error
}
