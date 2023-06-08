package InternalRepository

import (
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingCCTVRepository interface {
	FindBySn(sn string) entity.TbStagingCctv
	InsertWithTx(data entity.TbStagingCctv, tx *gorm.DB) error
}

type StagingStagingCCTVRepository struct {
	connection *gorm.DB
}

func NewStagingCCTVRepository(conn *gorm.DB) StagingCCTVRepository {
	return &StagingStagingCCTVRepository{
		connection: conn,
	}
}

func (db *StagingStagingCCTVRepository) FindBySn(sn string) entity.TbStagingCctv {
	var data entity.TbStagingCctv
	db.connection.Debug().Model(&data).Where("sn", sn).Find(&data)
	return data
}

func (db *StagingStagingCCTVRepository) InsertWithTx(data entity.TbStagingCctv, tx *gorm.DB) error {
	return tx.Create(&data).Error
}
