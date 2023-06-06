package PrestagingCRMRequest

type PostPrestaging struct {
	IdUploader              string `json:"idUploader" form:"idUploader" binding:"required"`
	Uploader                string `json:"uploader" form:"uploader" binding:"required"`
	Sn                      string `json:"sn" form:"sn" binding:"required"`
	ProjectName             string `json:"projectName" form:"projectName" binding:"required"`
	FotoMesinCrmFull        string `json:"fotoMesinCrmFull" form:"fotoMesinCrmFull" binding:"required"`
	FotoSnMesinCrm          string `json:"fotoSnMesinCrm" form:"fotoSnMesinCrm" binding:"required"`
	FotoCameraAtas          string `json:"fotoCameraAtas" form:"fotoCameraAtas" binding:"required"`
	FotoCameraCashOut       string `json:"fotoCameraCashOut" form:"fotoCameraCashOut" binding:"required"`
	FotoSystemInformationCu string `json:"fotoSystemInformationCu" form:"fotoSystemInformationCu" binding:"required"`
	FotoKapasitasHardisk    string `json:"fotoKapasitasHardisk" form:"fotoKapasitasHardisk" binding:"required"`
	FotoKunciCrm            string `json:"fotoKunciCrm" form:"fotoKunciCrm" binding:"required"`
	FotoClr                 string `json:"fotoClr" form:"fotoClr" binding:"required"`
	FotoPortPc              string `json:"fotoPortPc" form:"fotoPortPc" binding:"required"`
	FotoStikerBit           string `json:"fotoStikerBit" form:"fotoStikerBit" binding:"required"`
	FotoStrukErrorLogTest   string `json:"fotoStrukErrorLogTest" form:"fotoStrukErrorLogTest" binding:"required"`
	FotoCeklist             string `json:"fotoCeklist" form:"fotoCeklist" binding:"required"`
	SnCpu                   string `json:"snCpu" form:"snCpu" binding:"required"`
	SnClr                   string `json:"snClr" form:"snClr" binding:"required"`
	SnReceiptPrinter        string `json:"snReceiptPrinter" form:"snReceiptPrinter" binding:"required"`
	SnUr                    string `json:"snUr" form:"snUr" binding:"required"`
	SnBv                    string `json:"snBv" form:"snBv" binding:"required"`
	SnMonitor               string `json:"snMonitor" form:"snMonitor" binding:"required"`
	StatusDeadPixelMonitor  string `json:"statusDeadPixelMonitor" form:"statusDeadPixelMonitor" binding:"required"`
	Brand                   string `json:"brand" form:"brand" binding:"required"`
	TextNotes               string `json:"textNotes" form:"textNotes" binding:"required"`
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
	IdVerifikator           string `json:"idVerifikator" form:"idVerifikator" binding:"required"`
	Verifikator             string `json:"verifikator" form:"verifikator" binding:"required"`
	Sn                      string `json:"sn" form:"sn" binding:"required"`
	FotoMesinCrmFull        string `json:"fotoMesinCrmFull" form:"fotoMesinCrmFull" `
	FotoSnMesinCrm          string `json:"fotoSnMesinCrm" form:"fotoSnMesinCrm" `
	FotoCameraAtas          string `json:"fotoCameraAtas" form:"fotoCameraAtas" `
	FotoCameraCashOut       string `json:"fotoCameraCashOut" form:"fotoCameraCashOut" `
	FotoSystemInformationCu string `json:"fotoSystemInformationCu" form:"fotoSystemInformationCu" `
	FotoKapasitasHardisk    string `json:"fotoKapasitasHardisk" form:"fotoKapasitasHardisk" `
	FotoKunciCrm            string `json:"fotoKunciCrm" form:"fotoKunciCrm" `
	FotoClr                 string `json:"fotoClr" form:"fotoClr" `
	FotoPortPc              string `json:"fotoPortPc" form:"fotoPortPc" `
	FotoStikerBit           string `json:"fotoStikerBit" form:"fotoStikerBit" `
	FotoStrukErrorLogTest   string `json:"fotoStrukErrorLogTest" form:"fotoStrukErrorLogTest" `
	FotoCeklist             string `json:"fotoCeklist" form:"fotoCeklist" `
}
