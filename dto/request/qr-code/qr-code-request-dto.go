package QrCodeRequest

type GenerateQr struct {
	Sn string `json:"sn" binding:"required"`
}
