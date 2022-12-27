package repository

import (
	"go-api/dto/response/kprResponseDTO"
	"go-api/entity"

	"gorm.io/gorm"
)

type KPRRepository interface {
	Insert(data entity.TblPengajuanKprBySales) entity.TblPengajuanKprBySales
	Delete(id string) entity.TblPengajuanKprBySales
	KPRAllUser() []kprResponseDTO.PengajuanAllUser
}

type kprConnection struct {
	connection *gorm.DB
}

func NewKPRRepository(conn *gorm.DB) KPRRepository {
	return &kprConnection{
		connection: conn,
	}
}

func (db *kprConnection) Insert(data entity.TblPengajuanKprBySales) entity.TblPengajuanKprBySales {
	db.connection.Save(&data)
	return data
}

func (db *kprConnection) Delete(id string) entity.TblPengajuanKprBySales {
	var data entity.TblPengajuanKprBySales
	db.connection.Model(&entity.TblPengajuanKprBySales{}).Where("id", id).Update("status", "on_deleted").Find(&data)
	return data
}

func (db *kprConnection) KPRAllUser() []kprResponseDTO.PengajuanAllUser {
	var result []kprResponseDTO.PengajuanAllUser
	var table1, table2 []kprResponseDTO.PengajuanAllUser

	db.connection.Raw("SELECT email from tbl_pengajuan_kpr").Find(&table1)

	db.connection.Raw("SELECT tc.email from tbl_pengajuan_kpr_by_sales tpkbs join tbl_customer tc on tc.id = tpkbs.customer_id").Find(&table2)

	result = append(result, table1...)
	result = append(result, table2...)

	return result
}
