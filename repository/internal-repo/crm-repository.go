package InternalRepository

import (
	CRMResponseDTO "go-api/dto/response/crm"

	"gorm.io/gorm"
)

type CRMRepository interface {
	GetAll() []CRMResponseDTO.Response
}

type crmConnection struct {
	connection *gorm.DB
}

func NewCRMRepository(conn *gorm.DB) CRMRepository {
	return &crmConnection{
		connection: conn,
	}
}

func (db *crmConnection) GetAll() []CRMResponseDTO.Response {
	var result []CRMResponseDTO.Response
	db.connection.Raw("SELECT * FROM tb_master_tid").Find(&result)
	return result
}
