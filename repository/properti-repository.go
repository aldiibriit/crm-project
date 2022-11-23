package repository

import (
	propertiDTO "go-api/dto/properti"
	"go-api/dto/response/landingPageResponseDTO"

	"gorm.io/gorm"
)

type PropertiRepository interface {
	AdvancedFilter(sqlStr string) interface{}
	AdvancedFilter2(trxId string) []propertiDTO.TblImageProperti
	FindNearby() []landingPageResponseDTO.DetailPropertiDtoRes
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

func (db *propertiConnection) FindNearby() []landingPageResponseDTO.DetailPropertiDtoRes {
	var data []landingPageResponseDTO.DetailPropertiDtoRes
	db.connection.Raw("SELECT * from (SELECT tp2.*,tp2.id as projectId, haversine_km(-6.200000,106.816666, al.latitude, al.longitude) AS distance FROM tbl_project p  " +
		"INNER JOIN tbl_alamat_properti al ON p.alamat_properti_id = al.id  " +
		"INNER JOIN tbl_properti tp ON tp.project_id = p.id " +
		"JOIN tbl_project tp2 on tp2.id = tp.project_id " +
		"WHERE (tp.status != 'sold') AND (tp.status != 'deleted') AND (tp.stock_units > 0) AND (p.status = 'published')  " +
		"group by p.id  " +
		") AS properti " +
		"WHERE distance <= 200 ORDER BY distance ASC limit 10").Find(&data)

	return data
}
