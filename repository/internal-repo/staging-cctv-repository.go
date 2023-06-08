package InternalRepository

import (
	"database/sql"
	"go-api/entity"

	"gorm.io/gorm"
)

type StagingCCTVRepository interface {
	FindBySn(sn string) entity.TbStagingCctv
	InsertWithTx(data entity.TbStagingCctv, tx *gorm.DB) error
	UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error
	FindAllSubmittedData() []entity.TbStagingCctv
	FindRejectedData(idUploader string) []entity.TbStagingCctv
}

type stagingCCTVRepository struct {
	connection *gorm.DB
}

func NewStagingCCTVRepository(conn *gorm.DB) StagingCCTVRepository {
	return &stagingCCTVRepository{
		connection: conn,
	}
}

func (db *stagingCCTVRepository) FindBySn(sn string) entity.TbStagingCctv {
	var data entity.TbStagingCctv
	db.connection.Debug().Model(&data).Where("sn", sn).Find(&data)
	return data
}

func (db *stagingCCTVRepository) InsertWithTx(data entity.TbStagingCctv, tx *gorm.DB) error {
	return tx.Create(&data).Error
}

func (db *stagingCCTVRepository) UpdateWithTx(data map[string]interface{}, tx *gorm.DB) error {
	return tx.Debug().Model(&entity.TbStagingCctv{}).Where("sn", data["sn"].(string)).Updates(&data).Error
}

func (db *stagingCCTVRepository) FindAllSubmittedData() []entity.TbStagingCctv {
	var result []entity.TbStagingCctv
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus in ("TASK SUBMITTED","TASK REUPLOADED")`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_staging_cctv tpc where sn in (` + concatedSn.String + `)`).Find(&result)
	return result
}

func (db *stagingCCTVRepository) FindRejectedData(idUploader string) []entity.TbStagingCctv {
	var result []entity.TbStagingCctv
	var concatedSn sql.NullString
	db.connection.Raw(`select group_concat(cast(sn as char))as concatedSn  from (
		select a.sn,(select status_description from tb_log_activity tla2 where tla2.sn = a.sn order by created_at desc limit 1)as lastStatus from (
		select distinct sn from tb_log_activity tla 
	)a
	)b where b.lastStatus = "TASK REJECTED"`).Scan(&concatedSn)

	db.connection.Raw(`select * from tb_staging_cctv tpc where sn in (` + concatedSn.String + `) and id_uploader = '` + idUploader + `'`).Find(&result)
	return result
}
