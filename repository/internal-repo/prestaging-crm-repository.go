package InternalRepository

import (
	"database/sql"
	"go-api/entity"

	"gorm.io/gorm"
)

type PrestagingCRMRepository interface {
	Insert(data entity.TbPrestagingCrm) error
	InsertWithTx(data entity.TbPrestagingCrm, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
	FindAllSubmittedData() []entity.TbPrestagingCrm
	FindBySn(sn string) entity.TbPrestagingCrm
	FindRejectedData(idUploader string) []entity.TbPrestagingCrm
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

func (db *prestagingCRMConnection) FindAllSubmittedData() []entity.TbPrestagingCrm {
	var result []entity.TbPrestagingCrm
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus in ("TASK SUBMITTED","TASK REUPLOADED")`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_prestaging_crm tpc where sn in (` + concatedSn.String + `)`).Find(&result)
	return result
}

func (db *prestagingCRMConnection) FindBySn(sn string) entity.TbPrestagingCrm {
	var result entity.TbPrestagingCrm
	db.connection.Model(entity.TbPrestagingCrm{}).Where("sn", sn).Find(&result)
	return result
}

func (db *prestagingCRMConnection) FindRejectedData(idUploader string) []entity.TbPrestagingCrm {
	var result []entity.TbPrestagingCrm
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus = "TASK REJECTED"`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_prestaging_crm tpc where sn in (` + concatedSn.String + `) and id_uploader = '` + idUploader + `'`).Find(&result)
	return result
}
