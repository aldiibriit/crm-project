package entity

type TbPrestagingCctv struct {
	ID                                  uint64 `json:"id" gorm:"primary_key:auto_increment" `
	IdUploader                          string `json:"idUploader" `
	Uploader                            string `json:"uploader"`
	Sn                                  string `json:"sn"`
	ProjectName                         string `json:"projectName"`
	FotoNvrDanCamera                    string `json:"fotoNvrDanCamera"`
	FotoSnNvr                           string `json:"fotoSnNvr"`
	FotoSnCamera                        string `json:"fotoSnCamera"`
	FotoKelengkapanNvr                  string `json:"fotoKelengkapanNvr"`
	FotoKelengkapanCamera               string `json:"fotoKelengkapanCamera"`
	FotoInstallasiInputDanOutputPortNvr string `json:"fotoInstallasiInputDanOutputPortNvr"`
	FotoTampilanCamera                  string `json:"fotoTampilanCamera"`
	FotoSettingResolusi                 string `json:"fotoSettingResolusi"`
	FotoSettingMotion                   string `json:"fotoSettingMotion"`
	FotoStorage                         string `json:"fotoStorage"`
	FotoCheckBackup                     string `json:"fotoCheckBackup"`
	FotoStikerBit                       string `json:"fotoStikerBit"`
	SnNvr                               string `json:"snNvr"`
	SnCamera                            string `json:"snCamera"`
	UserNewNvr                          string `json:"userNewNvr"`
	PasswordNewNvr                      string `json:"passwordNewNvr"`
	StatusBarang                        string `json:"statusBarang"`
	Brand                               string `json:"brand"`
	TextNotes                           string `json:"textNotes"`
}
