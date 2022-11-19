package repository

import (
	propertiDTO "go-api/dto/properti"

	"gorm.io/gorm"
)

type PropertiRepository interface {
	AdvancedFilter(sqlStr string) interface{}
	AdvancedFilter2(trxId string) []propertiDTO.TblImageProperti
}

type propertiConnection struct {
	connection *gorm.DB
}

// NewPropertiRepository creates an instance BookRepository
func NewPropertiRepository(dbConn *gorm.DB) PropertiRepository {
	return &propertiConnection{
		connection: dbConn,
	}
}

func (db *propertiConnection) AdvancedFilter(sqlStr string) interface{} {
	var result []propertiDTO.AdvancedFilter
	db.connection.Raw(sqlStr).Preload("ImageProperti").Preload("ImageProject").Find(&result)
	return result
}

func (db *propertiConnection) AdvancedFilter2(trxId string) []propertiDTO.TblImageProperti {
	var imageProject []propertiDTO.TblImageProperti
	db.connection.Raw(`SELECT * FROM tbl_image_properti WHERE trx_id = '` + trxId + `'`).Find(&imageProject)
	return imageProject
}
