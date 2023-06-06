package entity

type TbPrestagingCrm struct {
	ID                      uint64 `json:"id" gorm:"primary_key:auto_increment" `
	IdUploader              string `json:"idUploader" `
	Uploader                string `json:"uploader"`
	Sn                      string `json:"sn"`
	ProjectName             string `json:"projectName"`
	FotoMesinCrmFull        string `json:"fotoMesinCrmFull"`
	FotoSnMesinCrm          string `json:"fotoSnMesinCrm"`
	FotoCameraAtas          string `json:"fotoCameraAtas"`
	FotoCameraCashOut       string `json:"fotoCameraCashOut"`
	FotoSystemInformationCu string `json:"fotoSystemInformationCu"`
	FotoKapasitasHardisk    string `json:"fotoKapasitasHardisk"`
	FotoKunciCrm            string `json:"fotoKunciCrm"`
	FotoClr                 string `json:"fotoClr"`
	FotoPortPc              string `json:"fotoPortPc"`
	FotoStikerBit           string `json:"fotoStikerBit"`
	FotoStrukErrorLogTest   string `json:"fotoStrukErrorLogTest"`
	FotoCeklist             string `json:"fotoCeklist"`
	SnCpu                   string `json:"snCpu"`
	SnClr                   string `json:"snClr"`
	SnReceiptPrinter        string `json:"snReceiptPrinter"`
	SnUr                    string `json:"snUr"`
	SnBv                    string `json:"snBv"`
	SnMonitor               string `json:"snMonitor"`
	StatusDeadPixelMonitor  string `json:"statusDeadPixelMonitor"`
	Brand                   string `json:"brand"`
	TextNotes               string `json:"textNotes"`
}
