package CCTVResponseDTO

type Cctv struct {
	SnCctv      string `json:"snCctv" gorm:"column:sn_cctv"`
	Type        string `json:"type"`
	NoSpk       string `json:"noSpk" gorm:"column:no_spk"`
	ProjectName string `json:"projectName" gorm:"column:project_name"`
	Year        string `json:"year"`
}
