package repository

import (
	propertiDTO "go-api/dto/properti"

	"gorm.io/gorm"
)

type ImagePropertiRepository interface {
	FindByTrxIdOrderByImageNameAsc(idTableProperti string) []propertiDTO.TblImageProperti
	FindByTrxIdOrderByImageNameAsc2(trxId string) []propertiDTO.TblImageProperti
	GetEncryptedImageProjectID(projectId int) string
}

type imagePropertiConnection struct {
	connection *gorm.DB
}

func NewImagePropertiRepository(conn *gorm.DB) ImagePropertiRepository {
	return &imagePropertiConnection{
		connection: conn,
	}
}

func (db *imagePropertiConnection) FindByTrxIdOrderByImageNameAsc(idTableProperti string) []propertiDTO.TblImageProperti {
	var result []propertiDTO.TblImageProperti
	db.connection.Raw(`select tip.* from tbl_properti tp
	join tbl_media_properti tmp on tp.media_properti_id = tmp.id
	join tbl_image_properti tip on tip.trx_id = tmp.image_properti_id
	where tip.trx_id = tmp.image_properti_id and tp.id = ?`, idTableProperti).Find(&result)
	return result
}

func (db *imagePropertiConnection) FindByTrxIdOrderByImageNameAsc2(trxId string) []propertiDTO.TblImageProperti {
	var result []propertiDTO.TblImageProperti
	db.connection.Raw(`select * from tbl_image_properti where trx_id = '` + trxId + `'`).Find(&result)
	return result
}

func (db *imagePropertiConnection) GetEncryptedImageProjectID(projectId int) string {
	var encryptedData string
	db.connection.Raw(`select tmp.image_id as encryptedData from tbl_project tp 
	join tbl_media_project tmp on tp.media_project_id = tmp.id 
	where tp.id = ?`, projectId).Take(&encryptedData)
	return encryptedData
}
