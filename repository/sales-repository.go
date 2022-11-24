package repository

import (
	"go-api/dto/response/salesResponseDTO"
	"go-api/entity"

	"gorm.io/gorm"
)

type SalesRepository interface {
	InsertRelation(data entity.TblSales) error
	FindByEmailDeveloper(emailDeveloper string) []salesResponseDTO.MISDeveloper
	MISSuperAdmin() []salesResponseDTO.MISSuperAdmin
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

func (db *salesConnection) MISSuperAdmin() []salesResponseDTO.MISSuperAdmin {
	var data []salesResponseDTO.MISSuperAdmin
	db.connection.Raw(`SELECT ts.sales_email ,ts.sales_name,(select json_extract(metadata,'$.name') from tbl_user tu2 where tu2.email = ts.developer_email)as metadata,tp.status,tp.jenis_properti,tp.tipe_properti
	FROM tbl_project tp
	JOIN tbl_sales ts on tp.email = ts.developer_email 
	JOIN tbl_user tu on tu.email = ts.sales_email 
	order by ts.sales_email 
	`).Find(&data)
	return data
}
