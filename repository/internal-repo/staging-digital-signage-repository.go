package InternalRepository

import (
	"database/sql"
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingDigitalSignageRepository interface {
	FindBySn(sn string) entity.TbStagingDigitalSignage
	InsertWithTx(data entity.TbStagingDigitalSignage, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
	FindAllSubmittedData() []entity.TbStagingDigitalSignage
	FindRejectedData(idUploader string) []entity.TbStagingDigitalSignage
}

type stagingDigitalSignageConnection struct {
	connection *gorm.DB
}

func NewStagingDigitalSignageRepository(conn *gorm.DB) StagingDigitalSignageRepository {
	return &stagingDigitalSignageConnection{
		connection: conn,
	}
}

func (db *stagingDigitalSignageConnection) FindBySn(sn string) entity.TbStagingDigitalSignage {
	var data entity.TbStagingDigitalSignage
	db.connection.Debug().Model(&data).Where("sn", sn).Find(&data)
	return data
}

func (db *stagingDigitalSignageConnection) InsertWithTx(data entity.TbStagingDigitalSignage, tx *gorm.DB) error {
	return tx.Create(&data).Error
}

func (db *stagingDigitalSignageConnection) UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error {
	return tx.Debug().Model(&entity.TbStagingDigitalSignage{}).Where("sn", data["sn"].(string)).Updates(&data).Error
}

func (db *stagingDigitalSignageConnection) FindAllSubmittedData() []entity.TbStagingDigitalSignage {
	var result []entity.TbStagingDigitalSignage
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus in ("TASK SUBMITTED","TASK REUPLOADED")`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_staging_digital_signage tpc where sn in (` + concatedSn.String + `)`).Find(&result)
	return result
}

func (db *stagingDigitalSignageConnection) FindRejectedData(idUploader string) []entity.TbStagingDigitalSignage {
	var result []entity.TbStagingDigitalSignage
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus = "TASK REJECTED"`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_staging_digital_signage tpc where sn in (` + concatedSn.String + `) and id_uploader = '` + idUploader + `'`).Find(&result)
	return result
}
