package InternalRepository

import (
	"database/sql"
	"go-api/entity"

	"gorm.io/gorm"
)

type PrestagingUPSRepository interface {
	Insert(data entity.TbPrestagingUps) error
	InsertWithTx(data entity.TbPrestagingUps, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
	FindAllSubmittedData() []entity.TbPrestagingUps
	FindBySn(sn string) entity.TbPrestagingUps
	FindRejectedData(idUploader string) []entity.TbPrestagingUps
}

type prestagingUPSConnection struct {
	connection *gorm.DB
}

func NewPrestagingUPSRepository(conn *gorm.DB) PrestagingUPSRepository {
	return &prestagingUPSConnection{
		connection: conn,
	}
}

func (db *prestagingUPSConnection) Insert(data entity.TbPrestagingUps) error {
	return db.connection.Debug().Create(&data).Error
}

func (db *prestagingUPSConnection) InsertWithTx(data entity.TbPrestagingUps, tx *gorm.DB) error {
	return tx.Debug().Create(&data).Error
}

func (db *prestagingUPSConnection) UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error {
	return tx.Debug().Model(&entity.TbPrestagingUps{}).Where("sn", data["sn"].(string)).Updates(&data).Error
}

func (db *prestagingUPSConnection) FindAllSubmittedData() []entity.TbPrestagingUps {
	var result []entity.TbPrestagingUps
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus in ("TASK SUBMITTED","TASK REUPLOADED")`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_prestaging_ups tpc where sn in (` + concatedSn.String + `)`).Find(&result)
	return result
}

func (db *prestagingUPSConnection) FindBySn(sn string) entity.TbPrestagingUps {
	var result entity.TbPrestagingUps
	db.connection.Model(entity.TbPrestagingUps{}).Where("sn", sn).Find(&result)
	return result
}

func (db *prestagingUPSConnection) FindRejectedData(idUploader string) []entity.TbPrestagingUps {
	var result []entity.TbPrestagingUps
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus = "TASK REJECTED"`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_prestaging_ups tpc where sn in (` + concatedSn.String + `) and id_uploader = '` + idUploader + `'`).Find(&result)
	return result
}
