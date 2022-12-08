package KPRRequestDTO

type PengajuanKPRRequest struct {
	NIK                    string  `json:"nik" binding:"required"`
	SalesID                int     `json:"salesID" binding:"required"`
	Name                   string  `json:"name" binding:"required"`
	TanggalLahir           string  `json:"tanggalLahir" binding:"required"`
	TempatLahir            string  `json:"tempatLahir" binding:"required"`
	MaritalStatus          string  `json:"maritalStatus" binding:"required"`
	JenisKelamin           string  `json:"jenisKelamin" binding:"required"`
	MobileNo               string  `json:"mobileNo" binding:"required"`
	AlamatKTP              string  `json:"alamatKTP" binding:"required"`
	AlamatDomisili         string  `json:"alamatDomisili" binding:"required"`
	KancaTerdekat          string  `json:"kancaTerdekat" binding:"required"`
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
	RefNoPengajuanBrispot  string  `json:"refNoPengajuanBrispot" `
	Status                 string  `json:"status" `
	StatusCodeBrispot      string  `json:"statusCodeBrispot" `
	StatusDescBrispot      string  `json:"statusDescBrispot" `
	TenorFixedRate         int     `json:"tenorFixedRate"`
	UangMuka               string  `json:"uangMuka" `
	UkerCode               string  `json:"ukerCode" `
	UkerName               string  `json:"ukerName" `
	PropertiId             int     `json:"propertiID"`
	CustomerID             int     `json:"customerID" gorm:"column:customer_id"`
}
