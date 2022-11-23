package salesResponseDTO

type MISSuperAdmin struct {
	SalesName     string `json:"salesName" gorm:"column:sales_name"`
	Metadata      string `json:"metadata"`
	Status        string `json:"status"`
	JenisProperti string `json:"jenisProperti" gorm:"column:jenis_properti"`
	TipeProperti  string `json:"tipeProperti" gorm:"column:tipe_properti"`
}
