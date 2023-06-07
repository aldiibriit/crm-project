package entity

type TbPrestagingDigitalSignage struct {
	ID                     uint64 `json:"id" gorm:"primary_key:auto_increment" `
	IdUploader             string `json:"idUploader" `
	Uploader               string `json:"uploader"`
	Sn                     string `json:"sn"`
	ProjectName            string `json:"projectName"`
	FotoLedFull            string `json:"FotoLedFull"`
	FotoSnLed              string `json:"FotoSnLed"`
	FotoTestDeadPixelPutih string `json:"FotoTestDeadPixelPutih"`
	FotoTestDeadPixelBiru  string `json:"FotoTestDeadPixelBiru"`
	FotoTestDeadPixelMerah string `json:"FotoTestDeadPixelMerah"`
	FotoTestDeadPixelHijau string `json:"FotoTestDeadPixelHijau"`
	FotoTestDeadPixelHitam string `json:"FotoTestDeadPixelHitam"`
	FotoStikerBit          string `json:"fotoStikerBit"`
	FotoKelengkapan        string `json:"fotoKelengkapan"`
	SnLed                  string `json:"snLed"`
	StatusBarang           string `json:"statusBarang"`
	Brand                  string `json:"brand"`
	TextNotes              string `json:"textNotes"`
}
