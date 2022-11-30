package KPRRequestDTO

type PengajuanKPRRequest struct {
	NIK            string `json:"nik" binding:"required"`
	SalesID        int    `json:"salesID" binding:"required"`
	Name           string `json:"name" binding:"required"`
	TanggalLahir   string `json:"tanggalLahir" binding:"required"`
	TempatLahir    string `json:"tempatLahir" binding:"required"`
	MaritalStatus  string `json:"maritalStatus" binding:"required"`
	JenisKelamin   string `json:"jenisKelamin" binding:"required"`
	MobileNo       string `json:"mobileNo" binding:"required"`
	AlamatKTP      string `json:"alamatKTP" binding:"required"`
	AlamatDomisili string `json:"alamatDomisili" binding:"required"`
	KancaTerdekat  string `json:"kancaTerdekat" binding:"required"`
}
