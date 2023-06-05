package PrestagingCRMRequest

type PostPrestaging struct {
	IdUploader               string `json:"idUploader" form:"idUploader" binding:"required"`
	Uploader                 string `json:"uploader" form:"uploader" binding:"required"`
	Sn                       string `json:"sn" form:"sn" binding:"required"`
	ProjectName              string `json:"projectName" form:"projectName" binding:"required"`
	FotoSnCrm                string `json:"fotoSnCrm" form:"fotoSnCrm" binding:"required"`
	FotoLembarKelengkapanCrm string `json:"fotoLembarKelengkapanCrm" form:"fotoLembarKelengkapanCrm" binding:"required"`
	FotoCheckKameraAtas      string `json:"fotoCheckKameraAtas" form:"fotoCheckKameraAtas" binding:"required"`
	FotoCheckKameraSamping   string `json:"fotoCheckKameraSamping" form:"fotoCheckKameraSamping" binding:"required"`
	FotoKunciCrm             string `json:"fotoKunciCrm" form:"fotoKunciCrm" binding:"required"`
	FotoStrukErrorLog        string `json:"fotoStrukErrorLog" form:"fotoStrukErrorLog" binding:"required"`
	FotoContactlessReader    string `json:"fotoContactlessReader" form:"fotoContactlessReader" binding:"required"`
	FotoKomponenPc           string `json:"fotoKomponenPc" form:"fotoKomponenPc" binding:"required"`
	FotoStikerBriit          string `json:"fotoStikerBriit" form:"fotoStikerBriit" binding:"required"`
	TanggalPrestaging        string `json:"tanggalPrestaging" form:"tanggalPrestaging" binding:"required"`
	TextNotes                string `json:"textNotes" form:"textNotes" binding:"required"`
}

type ApprovePrestaging struct {
	Sn string `json:"sn" binding:"required"`
}

type RejectPrestaging struct {
	IdVerifikator            string `json:"idVerifikator" form:"idVerifikator" binding:"required"`
	Verifikator              string `json:"verifikator" form:"verifikator" binding:"required"`
	Sn                       string `json:"sn" form:"sn" binding:"required"`
	FotoSnCrm                string `json:"fotoSnCrm" form:"fotoSnCrm"`
	FotoLembarKelengkapanCrm string `json:"fotoLembarKelengkapanCrm" form:"fotoLembarKelengkapanCrm"`
	FotoCheckKameraAtas      string `json:"fotoCheckKameraAtas" form:"fotoCheckKameraAtas"`
	FotoCheckKameraSamping   string `json:"fotoCheckKameraSamping" form:"fotoCheckKameraSamping"`
	FotoKunciCrm             string `json:"fotoKunciCrm" form:"fotoKunciCrm"`
	FotoStrukErrorLog        string `json:"fotoStrukErrorLog" form:"fotoStrukErrorLog"`
	FotoContactlessReader    string `json:"fotoContactlessReader" form:"fotoContactlessReader"`
	FotoKomponenPc           string `json:"fotoKomponenPc" form:"fotoKomponenPc"`
	FotoStikerBriit          string `json:"fotoStikerBriit" form:"fotoStikerBriit"`
}
