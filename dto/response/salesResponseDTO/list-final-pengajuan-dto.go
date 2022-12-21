package salesResponseDTO

type ListFinalPengajuanDTO struct {
	SalesName     string `json:"salesName"`
	DeveloperName string `json:"developerName" gorm:"column:developerName"`
	CustomerName  string `json:"customerName" gorm:"column:customerName"`
	CreatedAt     string `json:"createdAt"`
	JenisProperti string `json:"jenisProperti"`
	TipeProperti  string `json:"tipeProperti"`
}
