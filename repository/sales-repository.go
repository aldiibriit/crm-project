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
	GetIDByRefCode(refCode string) int
	FindByEmail(email string) entity.TblSales
	DraftDetail(request salesRequestDTO.DraftDetailRequest) salesResponseDTO.DraftDetailDTO
	EditDraftDetail(request salesRequestDTO.EditDraftDetailRequestDTO) error
	ListFinalPengajuan(sqlstr string, sqlstrCount string) ([]salesResponseDTO.ListFinalPengajuanDTO, int64)
	FindByEmailCustomer(email string) entity.TblCustomer
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

func (db *salesConnection) FindByEmailCustomer(email string) entity.TblCustomer {
	var customer entity.TblCustomer
	db.connection.Raw("SELECT * from tbl_customer where email = ?", email).Find(&customer)
	return customer
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

func (db *salesConnection) GetIDByRefCode(refCode string) int {
	var id int
	db.connection.Raw(`SELECT id from tbl_sales where refferal_code = '` + refCode + `'`).Scan(&id)
	return id
}

func (db *salesConnection) FindByEmail(email string) entity.TblSales {
	var data entity.TblSales
	db.connection.Raw(`SELECT * FROM tbl_sales where sales_email ='` + email + `' `).Find(&data)
	return data
}

func (db *salesConnection) DraftDetail(request salesRequestDTO.DraftDetailRequest) salesResponseDTO.DraftDetailDTO {
	var result salesResponseDTO.DraftDetailDTO
	db.connection.Raw(`SELECT * from tbl_pengajuan_kpr_by_sales tpkbs join tbl_customer tc on tc.id = tpkbs.customer_id where tpkbs.id = ` + request.ID + ``).Find(&result)
	return result
}

func (db *salesConnection) ListFinalPengajuan(sqlStr string, sqlStrCount string) ([]salesResponseDTO.ListFinalPengajuanDTO, int64) {
	var result []salesResponseDTO.ListFinalPengajuanDTO
	var totalData int64

	db.connection.Raw(sqlStr).Find(&result)
	db.connection.Raw(sqlStrCount).Count(&totalData)
	return result, totalData
}

func (db *salesConnection) EditDraftDetail(request salesRequestDTO.EditDraftDetailRequestDTO) error {
	var currentEmail string

	db.connection.Debug().Raw(`SELECT email from tbl_customer where id = ` + request.ID + ``).Scan(&currentEmail)

	tx := db.connection.Begin()
	if err := tx.Model(&entity.TblCustomer{}).Where("id", request.ID).Updates(&entity.TblCustomer{Email: request.Email, NIK: request.NIK}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
