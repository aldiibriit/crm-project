package InternalRepository

import (
	"database/sql"
	"go-api/entity"

	"gorm.io/gorm"
)

type PrestagingCCTVRepository interface {
	Insert(data entity.TbPrestagingCctv) error
	InsertWithTx(data entity.TbPrestagingCctv, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
	FindAllSubmittedData() []entity.TbPrestagingCctv
	FindBySn(sn string) entity.TbPrestagingCctv
	FindRejectedData(idUploader string) []entity.TbPrestagingCctv
}

type prestagingCCTVConnection struct {
	connection *gorm.DB
}

func NewPrestagingCCTVRepository(conn *gorm.DB) PrestagingCCTVRepository {
	return &prestagingCCTVConnection{
		connection: conn,
	}
}

func (db *prestagingCCTVConnection) Insert(data entity.TbPrestagingCctv) error {
	return db.connection.Debug().Create(&data).Error
}

func (db *prestagingCCTVConnection) InsertWithTx(data entity.TbPrestagingCctv, tx *gorm.DB) error {
	return tx.Debug().Create(&data).Error
}

func (db *prestagingCCTVConnection) UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error {
	return tx.Debug().Model(&entity.TbPrestagingCctv{}).Where("sn", data["sn"].(string)).Updates(&data).Error
}

func (db *prestagingCCTVConnection) FindAllSubmittedData() []entity.TbPrestagingCctv {
	var result []entity.TbPrestagingCctv
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus in ("TASK SUBMITTED","TASK REUPLOADED")`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_prestaging_cctv tpc where sn in (` + concatedSn.String + `)`).Find(&result)
	return result
}

func (db *prestagingCCTVConnection) FindBySn(sn string) entity.TbPrestagingCctv {
	var result entity.TbPrestagingCctv
	db.connection.Model(entity.TbPrestagingCctv{}).Where("sn", sn).Find(&result)
	return result
}

func (db *prestagingCCTVConnection) FindRejectedData(idUploader string) []entity.TbPrestagingCctv {
	var result []entity.TbPrestagingCctv
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus = "TASK REJECTED"`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_prestaging_cctv tpc where sn in (` + concatedSn.String + `) and id_uploader = '` + idUploader + `'`).Find(&result)
	return result
}
