package KPRRequestDTO

type PengajuanKPRRequest struct {
	NIK                    string  `json:"nik"`
	SalesID                int     `json:"salesID"`
	Name                   string  `json:"name"`
	TanggalLahir           string  `json:"tanggalLahir"`
	TempatLahir            string  `json:"tempatLahir"`
	MaritalStatus          string  `json:"maritalStatus"`
	JenisKelamin           string  `json:"jenisKelamin"`
	MobileNo               string  `json:"mobileNo"`
	AlamatKTP              string  `json:"alamatKTP"`
	AlamatDomisili         string  `json:"alamatDomisili"`
	KancaTerdekat          string  `json:"kancaTerdekat"`
	ApprovalPlafondBrispot float64 `json:"approvalPlafondBrispot"`
	BiayaAdminNominal      float64 `json:"biayaAdminNominal"`
	BiayaAdminPercentage   float64 `json:"biayaAdminPercentage"`
	BiayaProvisiNominal    float64 `json:"biayaProvisiNominal"`
	BiayaProvisiPercentage float64 `json:"biayaProvisiPercentage"`
	Dob                    string  `json:"dob"`
	Pob                    string  `json:"pob"`
	DpNominal              string  `json:"dpNominal" `
	DpPersen               string  `json:"dpPersen" `
	FixedRate              float64 `json:"fixedRate"`
	FloatingRate           float64 `json:"floatingRate"`
	GimmickID              int     `json:"gimmickID"`
	HargaProperti          string  `json:"hargaProperti" `
	JumlahPinjaman         string  `json:"jumlahPinjaman" `
	Note                   string  `json:"note" `
	// RefNoPengajuanBrispot  string  `json:"refNoPengajuanBrispot" `
	// Status                 string  `json:"status" `
	// StatusCodeBrispot      string  `json:"statusCodeBrispot" `
	// StatusDescBrispot      string  `json:"statusDescBrispot" `
	TenorFixedRate int    `json:"tenorFixedRate"`
	UangMuka       string `json:"uangMuka" `
	UkerCode       string `json:"ukerCode" `
	UkerName       string `json:"ukerName" `
	PropertiId     int    `json:"propertiID"`
	CustomerID     int    `json:"customerID" gorm:"column:customer_id"`
}
