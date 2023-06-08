package StagingUPSRequest

type PostStaging struct {
	IdUploader                string `json:"idUploader" form:"idUploader" binding:"required"`
	Uploader                  string `json:"uploader" form:"uploader" binding:"required"`
	Sn                        string `json:"sn" form:"sn" binding:"required"`
	ProjectName               string `json:"projectName" form:"projectName" binding:"required"`
	FotoUpsFull               string `json:"fotoUpsFull" form:"fotoUpsFull" binding:"required"`
	FotoKelengkapanUps        string `json:"fotoKelengkapanUps" form:"fotoKelengkapanUps" binding:"required"`
	FotoStikerBitDanSucofindo string `json:"fotoStikerBitDanSucofindo" form:"fotoStikerBitDanSucofindo" binding:"required"`
	TextNotes                 string `json:"textNotes" form:"textNotes" binding:"required"`
}

type ApproveStaging struct {
	Sn string `json:"sn" binding:"required"`
}

type FindBySn struct {
	Sn string `json:"sn" binding:"required"`
}

type FindRejectedData struct {
	IdUploader string `json:"idUploader" binding:"required"`
}

type RejectStaging struct {
	IdVerifikator             string `json:"idVerifikator" form:"idVerifikator" binding:"required"`
	Verifikator               string `json:"verifikator" form:"verifikator" binding:"required"`
	Sn                        string `json:"sn" form:"sn" binding:"required"`
	FotoUpsFull               string `json:"fotoUpsFull" form:"fotoUpsFull"`
	FotoKelengkapanUps        string `json:"fotoKelengkapanUps" form:"fotoKelengkapanUps"`
	FotoStikerBitDanSucofindo string `json:"fotoStikerBitDanSucofindo" form:"fotoStikerBitDanSucofindo"`
}
