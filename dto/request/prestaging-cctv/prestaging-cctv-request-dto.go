package PrestagingCCTVRequest

type PostPrestaging struct {
	IdUploader                          string `json:"idUploader" form:"idUploader" binding:"required"`
	Uploader                            string `json:"uploader" form:"uploader" binding:"required"`
	Sn                                  string `json:"sn" form:"sn" binding:"required"`
	ProjectName                         string `json:"projectName" form:"projectName" binding:"required"`
	FotoNvrDanCamera                    string `json:"fotoNvrDanCamera" form:"fotoNvrDanCamera" binding:"required"`
	FotoSnNvr                           string `json:"fotoSnNvr" form:"fotoSnNvr" binding:"required"`
	FotoSnCamera                        string `json:"fotoSnCamera" form:"fotoSnCamera" binding:"required"`
	FotoKelengkapanNvr                  string `json:"fotoKelengkapanNvr" form:"fotoKelengkapanNvr" binding:"required"`
	FotoKelengkapanCamera               string `json:"fotoKelengkapanCamera" form:"fotoKelengkapanCamera" binding:"required"`
	FotoInstallasiInputDanOutputPortNvr string `json:"fotoInstallasiInputDanOutputPortNvr" form:"fotoInstallasiInputDanOutputPortNvr" binding:"required"`
	FotoTampilanCamera                  string `json:"fotoTampilanCamera" form:"fotoTampilanCamera" binding:"required"`
	FotoSettingResolusi                 string `json:"fotoSettingResolusi" form:"fotoSettingResolusi" binding:"required"`
	FotoSettingMotion                   string `json:"fotoSettingMotion" form:"fotoSettingMotion" binding:"required"`
	FotoStorage                         string `json:"fotoStorage" form:"fotoStorage" binding:"required"`
	FotoCheckBackup                     string `json:"fotoCheckBackup" form:"fotoCheckBackup" binding:"required"`
	FotoStikerBit                       string `json:"fotoStikerBit" form:"fotoStikerBit" binding:"required"`
	SnNvr                               string `json:"snNvr" form:"snNvr" binding:"required"`
	SnCamera                            string `json:"snCamera" form:"snCamera" binding:"required"`
	UserNewNvr                          string `json:"userNewNvr" form:"userNewNvr" binding:"required"`
	PasswordNewNvr                      string `json:"passwordNewNvr" form:"passwordNewNvr" binding:"required"`
	Brand                               string `json:"brand" form:"brand" binding:"required"`
	StatusBarang                        string `json:"statusBarang" form:"statusBarang" binding:"required"`
	TextNotes                           string `json:"textNotes" form:"textNotes" binding:"required"`
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
	IdVerifikator                       string `json:"idVerifikator" form:"idVerifikator" binding:"required"`
	Verifikator                         string `json:"verifikator" form:"verifikator" binding:"required"`
	Sn                                  string `json:"sn" form:"sn" binding:"required"`
	FotoNvrDanCamera                    string `json:"fotoNvrDanCamera" form:"fotoNvrDanCamera" binding:"required"`
	FotoSnNvr                           string `json:"fotoSnNvr" form:"fotoSnNvr" binding:"required"`
	FotoSnCamera                        string `json:"fotoSnCamera" form:"fotoSnCamera" binding:"required"`
	FotoKelengkapanNvr                  string `json:"fotoKelengkapanNvr" form:"fotoKelengkapanNvr" binding:"required"`
	FotoKelengkapanCamera               string `json:"fotoKelengkapanCamera" form:"fotoKelengkapanCamera" binding:"required"`
	FotoInstallasiInputDanOutputPortNvr string `json:"fotoInstallasiInputDanOutputPortNvr" form:"fotoInstallasiInputDanOutputPortNvr" binding:"required"`
	FotoTampilanCamera                  string `json:"fotoTampilanCamera" form:"fotoTampilanCamera" binding:"required"`
	FotoSettingResolusi                 string `json:"fotoSettingResolusi" form:"fotoSettingResolusi" binding:"required"`
	FotoSettingMotion                   string `json:"fotoSettingMotion" form:"fotoSettingMotion" binding:"required"`
	FotoStorage                         string `json:"fotoStorage" form:"fotoStorage" binding:"required"`
	FotoCheckBackup                     string `json:"fotoCheckBackup" form:"fotoCheckBackup" binding:"required"`
	FotoStikerBit                       string `json:"fotoStikerBit" form:"fotoStikerBit" binding:"required"`
}
