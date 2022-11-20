package entity

import "time"

type TblPengajuanKprBySales struct {
	ID                      int `gorm:"primary_key:auto_increment" json:"id"`
	ApprovalPlafondBrispot  float64
	BiayaAdminNominal       float64
	BiayaAdminPercentage    float64
	BiayaProvinisiNominal   float64
	BiayaProvinsiPerventage float64
	Dob                     string `gorm:"size:255"`
	Pob                     string `gorm:"size:255"`
	DpNominal               string `gorm:"size:255"`
	DpPersen                string `gorm:"size:255"`
	FixedRate               float64
	FloatingRate            float64
	GimmickID               int
	HargaProperti           string `gorm:"size:255"`
	JumlahPinjaman          string `gorm:"size:255"`
	Note                    string `gorm:"size:255"`
	RefNoPengajuanBrispot   string `gorm:"size:255"`
	Status                  string `gorm:"size:255"`
	StatusCodeBrispot       string `gorm:"size:255"`
	StatusDescBrispot       string `gorm:"size:255"`
	TenorFixedRate          int
	UangMuka                string `gorm:"size:255"`
	UkerCode                string `gorm:"size:255"`
	UkerName                string `gorm:"size:255"`
	PropertiId              int
	CustomerID              int `gorm:"column:customer_id`
	CreatedAt               time.Time
	ModifiedAt              time.Time
}
