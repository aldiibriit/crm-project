package InternalRepository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingCRMRepository interface {
	FindBySn(sn string) entity.TbStagingCrm
	InsertWithTx(data entity.TbStagingCrm, tx *gorm.DB) error
}

type stagingCRMConnection struct {
	connection *gorm.DB
}

func NewStagingCRMRepository(conn *gorm.DB) StagingCRMRepository {
	return &stagingCRMConnection{
		connection: conn,
	}
}

func (db *stagingCRMConnection) FindBySn(sn string) entity.TbStagingCrm {
	var data entity.TbStagingCrm
	db.connection.Debug().Model(&data).Where("sn", sn).Find(&data)
	return data
}

func (db *stagingCRMConnection) InsertWithTx(data entity.TbStagingCrm, tx *gorm.DB) error {
	return tx.Create(&data).Error
}
