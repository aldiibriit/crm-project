package InternalRepository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingUPSRepository interface {
	FindBySn(sn string) entity.TbStagingUps
	InsertWithTx(data entity.TbStagingUps, tx *gorm.DB) error
}

type stagingUPSConnection struct {
	connection *gorm.DB
}

func NewStagingUPSRepository(conn *gorm.DB) StagingUPSRepository {
	return &stagingUPSConnection{
		connection: conn,
	}
}

func (db *stagingUPSConnection) FindBySn(sn string) entity.TbStagingUps {
	var data entity.TbStagingUps
	db.connection.Debug().Model(&data).Where("sn", sn).Find(&data)
	return data
}

func (db *stagingUPSConnection) InsertWithTx(data entity.TbStagingUps, tx *gorm.DB) error {
	return tx.Create(&data).Error
}
