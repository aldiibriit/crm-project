package InternalRepository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingDigitalSignageRepository interface {
	FindBySn(sn string) entity.TbStagingDigitalSignage
	InsertWithTx(data entity.TbStagingDigitalSignage, tx *gorm.DB) error
}

type StagingDigitalSigangeConnection struct {
	connection *gorm.DB
}

func NewStagingDigitalSignageRepository(conn *gorm.DB) StagingDigitalSignageRepository {
	return &StagingDigitalSigangeConnection{
		connection: conn,
	}
}

func (db *StagingDigitalSigangeConnection) FindBySn(sn string) entity.TbStagingDigitalSignage {
	var data entity.TbStagingDigitalSignage
	db.connection.Debug().Model(&data).Where("sn", sn).Find(&data)
	return data
}

func (db *StagingDigitalSigangeConnection) InsertWithTx(data entity.TbStagingDigitalSignage, tx *gorm.DB) error {
	return tx.Create(&data).Error
}
