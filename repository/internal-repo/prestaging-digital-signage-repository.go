package InternalRepository

import (
	"database/sql"
	"go-api/entity"

	"gorm.io/gorm"
)

type PrestagingDigitalSignageRepository interface {
	Insert(data entity.TbPrestagingDigitalSignage) error
	InsertWithTx(data entity.TbPrestagingDigitalSignage, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
	FindAllSubmittedData() []entity.TbPrestagingDigitalSignage
	FindBySn(sn string) entity.TbPrestagingDigitalSignage
	FindRejectedData(idUploader string) []entity.TbPrestagingDigitalSignage
}

type prestagingDigitalSignageConnection struct {
	connection *gorm.DB
}

func NewPrestagingDigitalSignageRepository(conn *gorm.DB) PrestagingDigitalSignageRepository {
	return &prestagingDigitalSignageConnection{
		connection: conn,
	}
}

func (db *prestagingDigitalSignageConnection) Insert(data entity.TbPrestagingDigitalSignage) error {
	return db.connection.Debug().Create(&data).Error
}

func (db *prestagingDigitalSignageConnection) InsertWithTx(data entity.TbPrestagingDigitalSignage, tx *gorm.DB) error {
	return tx.Debug().Create(&data).Error
}

func (db *prestagingDigitalSignageConnection) UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error {
	return tx.Debug().Model(&entity.TbPrestagingDigitalSignage{}).Where("sn", data["sn"].(string)).Updates(&data).Error
}

func (db *prestagingDigitalSignageConnection) FindAllSubmittedData() []entity.TbPrestagingDigitalSignage {
	var result []entity.TbPrestagingDigitalSignage
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus in ("TASK SUBMITTED","TASK REUPLOADED")`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_prestaging_digital_signage tpc where sn in (` + concatedSn.String + `)`).Find(&result)
	return result
}

func (db *prestagingDigitalSignageConnection) FindBySn(sn string) entity.TbPrestagingDigitalSignage {
	var result entity.TbPrestagingDigitalSignage
	db.connection.Model(entity.TbPrestagingDigitalSignage{}).Where("sn", sn).Find(&result)
	return result
}

func (db *prestagingDigitalSignageConnection) FindRejectedData(idUploader string) []entity.TbPrestagingDigitalSignage {
	var result []entity.TbPrestagingDigitalSignage
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus = "TASK REJECTED"`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_prestaging_digital_signage tpc where sn in (` + concatedSn.String + `) and id_uploader = '` + idUploader + `'`).Find(&result)
	return result
}
