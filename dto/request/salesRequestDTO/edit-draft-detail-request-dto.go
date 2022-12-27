package salesRequestDTO

type EditDraftDetailRequestDTO struct {
	ID             string `json:"id" binding:"required"`
	NIK            string `json:"nik"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	MobileNo       string `json:"mobile_no"`
	MartialStatus  string `json:"martial_status"`
	AlamatDomisili string `json:"alamat_domisili"`
	AlamatKTP      string `json:"alamat_ktp"`
}
