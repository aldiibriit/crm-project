package PrestagingDigitalSignageRequest

type PostPrestaging struct {
	IdUploader             string `json:"idUploader" form:"idUploader" binding:"required"`
	Uploader               string `json:"uploader" form:"uploader" binding:"required"`
	Sn                     string `json:"sn" form:"sn" binding:"required"`
	ProjectName            string `json:"projectName" form:"projectName" binding:"required"`
	FotoLedFull            string `json:"fotoLedFull" form:"fotoLedFull" binding:"required"`
	FotoSnLed              string `json:"fotoSnLed" form:"fotoSnLed" binding:"required"`
	FotoTestDeadPixelPutih string `json:"fotoTestDeadPixelPutih" form:"fotoTestDeadPixelPutih" binding:"required"`
	FotoTestDeadPixelBiru  string `json:"fotoTestDeadPixelBiru" form:"fotoTestDeadPixelBiru" binding:"required"`
	FotoTestDeadPixelMerah string `json:"fotoTestDeadPixelMerah" form:"fotoTestDeadPixelMerah" binding:"required"`
	FotoTestDeadPixelHijau string `json:"fotoTestDeadPixelHijau" form:"fotoTestDeadPixelHijau" binding:"required"`
	FotoTestDeadPixelHitam string `json:"fotoTestDeadPixelHitam" form:"fotoTestDeadPixelPutih" binding:"required"`
	FotoKelengkapan        string `json:"fotoKelengkapan" form:"fotoKelengkapan" binding:"required"`
	FotoStikerBit          string `json:"fotoStikerBit" form:"fotoStikerBit" binding:"required"`
	SnLed                  string `json:"snLed" form:"snLed" binding:"required"`
	Brand                  string `json:"brand" form:"brand" binding:"required"`
	StatusBarang           string `json:"statusBarang" form:"statusBarang" binding:"required"`
	TextNotes              string `json:"textNotes" form:"textNotes" binding:"required"`
}

type ApprovePrestaging struct {
	Sn string `json:"sn" binding:"required"`
}

type FindBySn struct {
	Sn string `json:"sn" binding:"required"`
}

type FindRejectedData struct {
	IdUploader string `json:"idUploader" binding:"required"`
}

type RejectPrestaging struct {
	IdVerifikator          string `json:"idVerifikator" form:"idVerifikator" binding:"required"`
	Verifikator            string `json:"verifikator" form:"verifikator" binding:"required"`
	Sn                     string `json:"sn" form:"sn" binding:"required"`
	FotoLedFull            string `json:"fotoLedFull" form:"fotoLedFull" binding:"required"`
	FotoSnLed              string `json:"fotoSnLed" form:"fotoSnLed" binding:"required"`
	FotoTestDeadPixelPutih string `json:"fotoTestDeadPixelPutih" form:"fotoTestDeadPixelPutih" binding:"required"`
	FotoTestDeadPixelBiru  string `json:"fotoTestDeadPixelBiru" form:"fotoTestDeadPixelBiru" binding:"required"`
	FotoTestDeadPixelMerah string `json:"fotoTestDeadPixelMerah" form:"fotoTestDeadPixelMerah" binding:"required"`
	FotoTestDeadPixelHijau string `json:"fotoTestDeadPixelHijau" form:"fotoTestDeadPixelHijau" binding:"required"`
	FotoTestDeadPixelHitam string `json:"fotoTestDeadPixelHitam" form:"fotoTestDeadPixelPutih" binding:"required"`
	FotoKelengkapan        string `json:"fotoKelengkapan" form:"fotoKelengkapan" binding:"required"`
	FotoStikerBit          string `json:"fotoStikerBit" form:"fotoStikerBit" binding:"required"`
}
