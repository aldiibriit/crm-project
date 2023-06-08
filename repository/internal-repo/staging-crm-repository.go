package InternalRepository

import (
	"database/sql"
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingCRMRepository interface {
	FindBySn(sn string) entity.TbStagingCrm
	InsertWithTx(data entity.TbStagingCrm, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
	FindAllSubmittedData() []entity.TbStagingCrm
	FindRejectedData(idUploader string) []entity.TbStagingCrm
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

func (db *stagingCRMConnection) UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error {
	return tx.Debug().Model(&entity.TbStagingCrm{}).Where("sn", data["sn"].(string)).Updates(&data).Error
}

func (db *stagingCRMConnection) FindAllSubmittedData() []entity.TbStagingCrm {
	var result []entity.TbStagingCrm
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus in ("TASK SUBMITTED","TASK REUPLOADED")`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_staging_crm tpc where sn in (` + concatedSn.String + `)`).Find(&result)
	return result
}

func (db *stagingCRMConnection) FindRejectedData(idUploader string) []entity.TbStagingCrm {
	var result []entity.TbStagingCrm
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus = "TASK REJECTED"`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_staging_crm tpc where sn in (` + concatedSn.String + `) and id_uploader = '` + idUploader + `'`).Find(&result)
	return result
}
