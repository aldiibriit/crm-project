package KPRRequestDTO

type PengajuanKPRRequest struct {
	Name                   string  `json:"name" binding:"required"`
	ReferralCode           string  `json:"referralCode" binding:"required"`
	MobileNo               string  `json:"mobileNo" binding:"required"`
	PropertiId             int     `json:"propertiID" binding:"required"`
	Periode                string  `json:"periode"`
	NIK                    string  `json:"nik"`
	SalesEmail             string  `json:"salesEmail"`
	SalesID                int     `json:"salesID"`
	TanggalLahir           string  `json:"tanggalLahir"`
	TempatLahir            string  `json:"tempatLahir"`
	Email                  string  `json:"email"`
	MaritalStatus          string  `json:"maritalStatus"`
	JenisKelamin           string  `json:"jenisKelamin"`
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
	TenorFixedRate         int     `json:"tenorFixedRate"`
	CustomerID             int     `json:"customerID" gorm:"column:customer_id"`
	UangMuka               string  `json:"uangMuka" `
	UkerCode               string  `json:"ukerCode" `
	UkerName               string  `json:"ukerName" `
	UTJ                    string  `json:"utj"`
}
