package InternalRepository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingLiveCRMRepository interface {
	FindBySn(sn string) entity.TbStagingLiveCrm
	InsertWithTx(data entity.TbStagingLiveCrm, tx *gorm.DB) error
}

type stagingLiveCRMConnection struct {
	connection *gorm.DB
}

func NewStagingLiveCRMRepository(conn *gorm.DB) StagingLiveCRMRepository {
	return &stagingLiveCRMConnection{
		connection: conn,
	}
}

func (db *stagingLiveCRMConnection) FindBySn(sn string) entity.TbStagingLiveCrm {
	var data entity.TbStagingLiveCrm
	db.connection.Debug().Model(&data).Where("sn", sn).Find(&data)
	return data
}

func (db *stagingLiveCRMConnection) InsertWithTx(data entity.TbStagingLiveCrm, tx *gorm.DB) error {
	return tx.Create(&data).Error
}
