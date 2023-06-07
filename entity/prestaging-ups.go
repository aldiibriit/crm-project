package entity

type TbPrestagingUps struct {
	ID                                uint64 `json:"id" gorm:"primary_key:auto_increment" `
	IdUploader                        string `json:"idUploader" `
	Uploader                          string `json:"uploader"`
	Sn                                string `json:"sn"`
	ProjectName                       string `json:"projectName"`
	FotoUpsFull                       string `json:"fotoUpsFull"`
	FotoSnUps                         string `json:"fotoSnUps"`
	FotoKapasitasBaterai              string `json:"fotoKapasitasBaterai"`
	FotoTampilanKelistrikanListrikOn  string `json:"fotoTampilanKelistrikanListrikOn"`
	FotoTampilanKelistrikanListrikOff string `json:"fotoTampilanKelistrikanListrikOff"`
	FotoKelengkapan                   string `json:"fotoKelengkapan"`
	FotoStikerBit                     string `json:"fotoStikerBit"`
	FotoCeklist                       string `json:"fotoCeklist"`
	StatusBarang                      string `json:"statusBarang"`
	Brand                             string `json:"brand"`
	TextNotes                         string `json:"textNotes"`
}
