package entity

type TbStagingDigitalSignage struct {
	ID                        uint64 `json:"id" gorm:"primary_key:auto_increment" `
	IdUploader                string `json:"idUploader"`
	Uploader                  string `json:"uploader"`
	Sn                        string `json:"sn"`
	ProjectName               string `json:"projectName"`
	FotoLed                   string `json:"fotoLed"`
	FotoSnLedDus              string `json:"fotoSnLedDus"`
	FotoStikerBitDanSucofindo string `json:"fotoStikerBitDanSucofindo"`
}
