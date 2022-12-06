package repository

import (
	"go-api/dto/request/salesRequestDTO"
	"go-api/dto/response/salesResponseDTO"
	"go-api/entity"

	"gorm.io/gorm"
)

type SalesRepository interface {
	InsertRelation(data entity.TblSales) error
	FindByEmailDeveloper(sqlStr string, sqlStrCount string, request salesRequestDTO.MISDeveloperRequestDTO) ([]salesResponseDTO.MISDeveloper, int)
	MISSuperAdmin(sqlStr string, sqlStrCount string, request salesRequestDTO.MISSuperAdminRequestDTO) ([]salesResponseDTO.MISSuperAdmin, int)
	ListProject(sqlStr string) ([]salesResponseDTO.ListProject, int64)
	RelationToImageProperti(trxId string) []salesResponseDTO.TblImageProperti
	DetailSalesByDeveloper(request salesRequestDTO.DetailSalesRequest) salesResponseDTO.MISDeveloper
	EditSalesByDeveloper(request salesRequestDTO.SalesEditRequestDTO) error
	DeleteSalesByDeveloper(request salesRequestDTO.SalesDeleteRequestDTO) error
}

type salesConnection struct {
	connection *gorm.DB
}

func NewSalesRepository(conn *gorm.DB) SalesRepository {
	return &salesConnection{
		connection: conn,
	}
}

func (db *salesConnection) InsertRelation(data entity.TblSales) error {
	err := db.connection.Save(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *salesConnection) FindByEmailDeveloper(sqlStr string, sqlStrCount string, request salesRequestDTO.MISDeveloperRequestDTO) ([]salesResponseDTO.MISDeveloper, int) {
	var result []salesResponseDTO.MISDeveloper
	var totalData int

	db.connection.Raw(sqlStr).Find(&result)

	db.connection.Raw(sqlStrCount).Find(&totalData)

	return result, totalData
}

func (db *salesConnection) MISSuperAdmin(sqlStr string, sqlStrCount string, request salesRequestDTO.MISSuperAdminRequestDTO) ([]salesResponseDTO.MISSuperAdmin, int) {
	var data []salesResponseDTO.MISSuperAdmin
	var totalData int

	db.connection.Raw(sqlStr).Find(&data)

	db.connection.Raw(sqlStrCount).Find(&totalData)

	return data, totalData
}

func (db *salesConnection) ListProject(sqlStr string) ([]salesResponseDTO.ListProject, int64) {
	var result []salesResponseDTO.ListProject
	var length int64
	db.connection.Raw(sqlStr).Find(&result)
	db.connection.Raw(sqlStr).Find(&result).Count(&length)
	return result, length
}

func (db *salesConnection) RelationToImageProperti(trxId string) []salesResponseDTO.TblImageProperti {
	var imageProject []salesResponseDTO.TblImageProperti
	db.connection.Raw(`SELECT * FROM tbl_image_properti WHERE trx_id = '` + trxId + `'`).Find(&imageProject)
	return imageProject
}

func (db *salesConnection) EditSalesByDeveloper(request salesRequestDTO.SalesEditRequestDTO) error {
	var currentEmail string

	db.connection.Raw(`SELECT email from tbl_user where id = ` + request.ID + ``).Scan(&currentEmail)

	tx := db.connection.Begin()
	if err := tx.Model(&entity.TblUser{}).Where("id", request.ID).Updates(&entity.TblUser{Email: request.Email, MobileNo: request.SalesPhone}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&entity.TblSales{}).Where("sales_email", currentEmail).Update("sales_email", request.Email).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (db *salesConnection) DeleteSalesByDeveloper(request salesRequestDTO.SalesDeleteRequestDTO) error {
	tx := db.connection.Begin()
	if err := tx.Model(&entity.TblUser{}).Where("id", request.ID).Update("status", "Deleted").Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (db *salesConnection) DetailSalesByDeveloper(request salesRequestDTO.DetailSalesRequest) salesResponseDTO.MISDeveloper {
	var result salesResponseDTO.MISDeveloper
	db.connection.Raw(`SELECT 
	tu.id,ts.developer_email,ts.sales_email,ts.refferal_code,ts.registered_by,ts.created_at,ts.modified_at,ts.sales_name,tu.mobile_no as salesPhone,tu.status
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE tu.id = ` + request.ID + ``).First(&result)
	return result
}
