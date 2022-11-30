package entity

import "time"

type TblCustomer struct {
	ID             int    `gorm:"primary_key:auto_increment" json:"id"`
	NIK            string `json:"nik" gorm:"column:nik;size:255"`
	Name           string `json:"name" gorm:"column:name;size:255"`
	Email          string `json:"email" gorm:"column:email;size:255"`
	MobileNo       string `json:"mobileNo" gorm:"column:mobile_no;size:255"`
	MaritalStatus  string `json:"maritalStatus" gorm:"column:marital_status;size:255"`
	AlamatDomisili string `json:"alamatDomisili" gorm:"column:alamat_domisili;type:longtext"`
	AlamatKTP      string `json:"alamatKTP" gorm:"column:alamat_ktp;type:longtext"`
	SalesID        int    `json:"salesID" gorm:"column:sales_id"`
	CreatedAt      time.Time
	ModifiedAt     time.Time
	// PengajuanKPRBySales TblPengajuanKprBySales
}
