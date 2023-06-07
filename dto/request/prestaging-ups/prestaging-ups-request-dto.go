package PrestagingUPSRequest

type PostPrestaging struct {
	IdUploader                        string `json:"idUploader" form:"idUploader" binding:"required"`
	Uploader                          string `json:"uploader" form:"uploader" binding:"required"`
	Sn                                string `json:"sn" form:"sn" binding:"required"`
	ProjectName                       string `json:"projectName" form:"projectName" binding:"required"`
	FotoUpsFull                       string `json:"fotoUpsFull" form:"fotoUpsFull" binding:"required"`
	FotoSnUps                         string `json:"fotoSnUps" form:"fotoSnUps" binding:"required"`
	FotoKapasitasBaterai              string `json:"fotoKapasitasBaterai" form:"fotoKapasitasBaterai" binding:"required"`
	FotoTampilanKelistrikanListrikOn  string `json:"fotoTampilanKelistrikanListrikOn" form:"fotoTampilanKelistrikanListrikOn" binding:"required"`
	FotoTampilanKelistrikanListrikOff string `json:"fotoTampilanKelistrikanListrikOff" form:"fotoTampilanKelistrikanListrikOff" binding:"required"`
	FotoKelengkapan                   string `json:"fotoKelengkapan" form:"fotoKelengkapan" binding:"required"`
	FotoStikerBit                     string `json:"fotoStikerBit" form:"fotoStikerBit" binding:"required"`
	FotoCeklist                       string `json:"fotoCeklist" form:"fotoCeklist" binding:"required"`
	FotoTampakBelakang                string `json:"fotoTampakBelakang" form:"fotoTampakBelakang" binding:"required"`
	StatusBarang                      string `json:"statusBarang" form:"statusBarang" binding:"required"`
	Brand                             string `json:"brand" form:"brand" binding:"required"`
	TextNotes                         string `json:"textNotes" form:"textNotes" binding:"required"`
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
	IdVerifikator                     string `json:"idVerifikator" form:"idVerifikator" binding:"required"`
	Verifikator                       string `json:"verifikator" form:"verifikator" binding:"required"`
	Sn                                string `json:"sn" form:"sn" binding:"required"`
	FotoUpsFull                       string `json:"fotoUpsFull" form:"fotoUpsFull" binding:"required"`
	FotoSnUps                         string `json:"fotoSnUps" form:"fotoSnUps" binding:"required"`
	FotoKapasitasBaterai              string `json:"fotoKapasitasBaterai" form:"fotoKapasitasBaterai" binding:"required"`
	FotoTampilanKelistrikanListrikOn  string `json:"fotoTampilanKelistrikanListrikOn" form:"fotoTampilanKelistrikanListrikOn" binding:"required"`
	FotoTampilanKelistrikanListrikOff string `json:"fotoTampilanKelistrikanListrikOff" form:"fotoTampilanKelistrikanListrikOff" binding:"required"`
	FotoTampakBelakang                string `json:"fotoTampakBelakang" form:"fotoTampakBelakang" binding:"required"`
	FotoKelengkapan                   string `json:"fotoKelengkapan" form:"fotoKelengkapan" binding:"required"`
	FotoStikerBit                     string `json:"fotoStikerBit" form:"fotoStikerBit" binding:"required"`
	FotoCeklist                       string `json:"fotoCeklist" form:"fotoCeklist" binding:"required"`
}
