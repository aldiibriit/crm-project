package repository

import (
	"go-api/dto/request/salesRequestDTO"
	"go-api/dto/response/salesResponseDTO"
	"go-api/entity"
	"strconv"

	"gorm.io/gorm"
)

type SalesRepository interface {
	InsertRelation(data entity.TblSales) error
	FindByEmailDeveloper(request salesRequestDTO.MISDeveloperRequestDTO) ([]salesResponseDTO.MISDeveloper, int)
	MISSuperAdmin(request salesRequestDTO.MISSuperAdminRequestDTO) ([]salesResponseDTO.MISSuperAdmin, int)
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

func (db *salesConnection) FindByEmailDeveloper(request salesRequestDTO.MISDeveloperRequestDTO) ([]salesResponseDTO.MISDeveloper, int) {
	var result []salesResponseDTO.MISDeveloper
	var totalData int

	db.connection.Raw(`SELECT 
	tu.id,ts.developer_email,ts.sales_email,ts.refferal_code,ts.registered_by,ts.created_at,ts.modified_at,ts.sales_name,tu.mobile_no as salesPhone
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = '` + request.EmailDeveloper + `' and ts.developer_email like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.refferal_code like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.registered_by like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.created_at like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.modified_at like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.sales_name like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and tu.mobile_no like '%` + request.Keyword + `%'
	limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
	`).Find(&result)

	db.connection.Raw(`SELECT 
	count(tu.id)
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = '` + request.EmailDeveloper + `' and ts.developer_email like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.refferal_code like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.registered_by like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.created_at like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.modified_at like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and ts.sales_name like '%` + request.Keyword + `%'
	or developer_email = '` + request.EmailDeveloper + `' and tu.mobile_no like '%` + request.Keyword + `%'
	`).Find(&totalData)

	return result, totalData
}

func (db *salesConnection) MISSuperAdmin(request salesRequestDTO.MISSuperAdminRequestDTO) ([]salesResponseDTO.MISSuperAdmin, int) {
	var data []salesResponseDTO.MISSuperAdmin
	var totalData int
	db.connection.Raw(`SELECT ts.sales_email ,ts.sales_name,(select json_extract(metadata,'$.name') from tbl_user tu2 where tu2.email = ts.developer_email)as metadata,tp.status,tp.jenis_properti,tp.tipe_properti
	FROM tbl_project tp
	JOIN tbl_sales ts on tp.email = ts.developer_email 
	JOIN tbl_user tu on tu.email = ts.sales_email 
	WHERE ts.sales_email like '%` + request.Keyword + `%' or ts.sales_name like '%` + request.Keyword + `%' or metadata like '%` + request.Keyword + `%' or tp.status like '%` + request.Keyword + `%' or tp.jenis_properti like '%` + request.Keyword + `%' or tp.tipe_properti like '%` + request.Keyword + `%'
	order by ts.sales_email limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
	`).Find(&data)

	db.connection.Raw(`SELECT count(ts.sales_email)
	FROM tbl_project tp
	JOIN tbl_sales ts on tp.email = ts.developer_email 
	JOIN tbl_user tu on tu.email = ts.sales_email 
	WHERE ts.sales_email like '%` + request.Keyword + `%' or ts.sales_name like '%` + request.Keyword + `%' or metadata like '%` + request.Keyword + `%' or tp.status like '%` + request.Keyword + `%' or tp.jenis_properti like '%` + request.Keyword + `%' or tp.tipe_properti like '%` + request.Keyword + `%'
	order by ts.sales_email limit ` + strconv.Itoa(request.Limit) + `
	`).Find(&totalData)

	return data, totalData
}
