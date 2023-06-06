package InternalRepository

import (
	"go-api/entity"

	LogActivityResponseDTO "go-api/dto/response/log-activity"

	"gorm.io/gorm"
)

type LogActivityRepository interface {
	Insert(data entity.TbLogActivity) error
	InsertWithTx(data entity.TbLogActivity, tx *gorm.DB) error
	FindTimeLineBySn(sn string) []LogActivityResponseDTO.DetailTimeline
}

type logActivityConnection struct {
	connection *gorm.DB
}

func NewLogActivityRepository(conn *gorm.DB) LogActivityRepository {
	return &logActivityConnection{
		connection: conn,
	}
}

func (db *logActivityConnection) Insert(data entity.TbLogActivity) error {
	return db.connection.Debug().Create(&data).Error
}

func (db *logActivityConnection) InsertWithTx(data entity.TbLogActivity, tx *gorm.DB) error {
	return tx.Debug().Create(&data).Error
}

func (db *logActivityConnection) FindTimeLineBySn(sn string) []LogActivityResponseDTO.DetailTimeline {
	var result []LogActivityResponseDTO.DetailTimeline
	db.connection.Raw(`select *,
	(select created_at from tb_log_activity tla2 where tla2.sn = a.sn and status_description = "TASK SUBMITTED" and a.category = tla2.category)as submittedAt,
		(select created_at from tb_log_activity tla2 where tla2.sn = a.sn and status_description = "TASK REJECTED" and a.category = tla2.category)as rejectedAt,
		(select created_at from tb_log_activity tla2 where tla2.sn = a.sn and status_description = "TASK REUPLOADED" and a.category = tla2.category)as reuploadedAt,
		(select created_at from tb_log_activity tla2 where tla2.sn = a.sn and status_description = "TASK APPROVED" and a.category = tla2.category)as approvedAt,
		category
	from (
		select distinct sn,category from tb_log_activity tla where sn = '` + sn + `' order by created_at asc
	)a`).Find(&result)
	return result
}
