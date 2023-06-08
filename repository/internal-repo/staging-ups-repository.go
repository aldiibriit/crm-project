package InternalRepository

import (
	"database/sql"
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingUPSRepository interface {
	FindBySn(sn string) entity.TbStagingUps
	InsertWithTx(data entity.TbStagingUps, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
	FindAllSubmittedData() []entity.TbStagingUps
	FindRejectedData(idUploader string) []entity.TbStagingUps
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

func (db *stagingUPSConnection) UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error {
	return tx.Debug().Model(&entity.TbStagingUps{}).Where("sn", data["sn"].(string)).Updates(&data).Error
}

func (db *stagingUPSConnection) FindAllSubmittedData() []entity.TbStagingUps {
	var result []entity.TbStagingUps
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus in ("TASK SUBMITTED","TASK REUPLOADED")`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_staging_ups tpc where sn in (` + concatedSn.String + `)`).Find(&result)
	return result
}

func (db *stagingUPSConnection) FindRejectedData(idUploader string) []entity.TbStagingUps {
	var result []entity.TbStagingUps
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus = "TASK REJECTED"`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_staging_ups tpc where sn in (` + concatedSn.String + `) and id_uploader = '` + idUploader + `'`).Find(&result)
	return result
}
