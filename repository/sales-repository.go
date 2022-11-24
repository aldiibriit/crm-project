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
	FindByEmailDeveloper(emailDeveloper string) []salesResponseDTO.MISDeveloper
	MISSuperAdmin(request salesRequestDTO.MISSuperAdminRequestDTO) []salesResponseDTO.MISSuperAdmin
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

func (db *salesConnection) FindByEmailDeveloper(emailDeveloper string) []salesResponseDTO.MISDeveloper {
	var result []salesResponseDTO.MISDeveloper
	db.connection.Raw(`SELECT 
	tu.id,ts.developer_email,ts.sales_email,ts.refferal_code,ts.registered_by,ts.created_at,ts.modified_at,ts.sales_name,tu.mobile_no as salesPhone
	FROM tbl_sales ts
	JOIN tbl_user tu ON tu.email = ts.sales_email
	WHERE developer_email = ?`, emailDeveloper).Find(&result)
	return result
}

func (db *salesConnection) MISSuperAdmin(request salesRequestDTO.MISSuperAdminRequestDTO) []salesResponseDTO.MISSuperAdmin {
	var data []salesResponseDTO.MISSuperAdmin
	db.connection.Raw(`SELECT ts.sales_email ,ts.sales_name,(select json_extract(metadata,'$.name') from tbl_user tu2 where tu2.email = ts.developer_email)as metadata,tp.status,tp.jenis_properti,tp.tipe_properti
	FROM tbl_project tp
	JOIN tbl_sales ts on tp.email = ts.developer_email 
	JOIN tbl_user tu on tu.email = ts.sales_email 
	WHERE ts.sales_email like '%` + request.Keyword + `%' or ts.sales_name like '%` + request.Keyword + `%' or metadata like '%` + request.Keyword + `%' or tp.status like '%` + request.Keyword + `%' or tp.jenis_properti like '%` + request.Keyword + `%' or tp.tipe_properti like '%` + request.Keyword + `%'
	order by ts.sales_email limit ` + strconv.Itoa(request.Limit) + ` offset ` + strconv.Itoa(request.Offset) + `
	`).Find(&data)
	return data
}
