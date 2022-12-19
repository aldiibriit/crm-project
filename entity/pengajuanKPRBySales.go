package entity

import "time"

type TblPengajuanKprBySales struct {
	ID                     int     `gorm:"primary_key:auto_increment" json:"id"`
	ApprovalPlafondBrispot float64 `json:"approvalPlafondBrispot"`
	BiayaAdminNominal      float64 `json:"biayaAdminNominal"`
	BiayaAdminPercentage   float64 `json:"biayaAdminPercentage"`
	BiayaProvisiNominal    float64 `json:"biayaProvisiNominal"`
	BiayaProvisiPercentage float64 `json:"biayaProvisiPercentage"`
	Dob                    string  `json:"dob" gorm:"size:255"`
	Pob                    string  `json:"pob" gorm:"size:255"`
	DpNominal              string  `json:"dpNominal" gorm:"size:255"`
	DpPersen               string  `json:"dpPersen" gorm:"size:255"`
	FixedRate              float64 `json:"fixedRate"`
	FloatingRate           float64 `json:"floatingRate"`
	GimmickID              int     `json:"gimmickID"`
	HargaProperti          string  `json:"hargaProperti" gorm:"size:255"`
	JumlahPinjaman         string  `json:"jumlahPinjaman" gorm:"size:255"`
	Note                   string  `json:"note" gorm:"size:255"`
	RefNoPengajuanBrispot  string  `json:"refNoPengajuanBrispot" gorm:"size:255"`
	Status                 string  `json:"status" gorm:"size:255"`
	StatusCodeBrispot      string  `json:"statusCodeBrispot" gorm:"size:255"`
	StatusDescBrispot      string  `json:"statusDescBrispot" gorm:"size:255"`
	TenorFixedRate         int     `json:"tenorFixedRate"`
	UangMuka               string  `json:"uangMuka" gorm:"size:255"`
	UkerCode               string  `json:"ukerCode" gorm:"size:255"`
	UkerName               string  `json:"ukerName" gorm:"size:255"`
	PropertiId             int     `json:"propertiID"`
	CustomerID             int     `json:"customerID" gorm:"column:customer_id"`
	UTJ                    string  `json:"utj" gorm:"column:utj"`
	UTJStatus              bool    `json:"utjStatus" gorm:"column:utj_status"`
	CreatedAt              time.Time
	ModifiedAt             time.Time
}
