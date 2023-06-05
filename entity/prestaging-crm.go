package entity

type TbPrestagingCrm struct {
	ID                       uint64 `json:"id" gorm:"primary_key:auto_increment" `
	IdUploader               string `json:"idUploader"`
	Uploader                 string `json:"uploader"`
	IdVerifikator            string `json:"idVerifikator"`
	Verifikator              string `json:"verifikator"`
	Sn                       string `json:"sn"`
	ProjectName              string `json:"projectName"`
	FotoSnCrm                string `json:"fotoSnCrm"`
	FotoLembarKelengkapanCrm string `json:"fotoLembarKelengkapanCrm"`
	FotoCheckKameraAtas      string `json:"fotoCheckKameraAtas"`
	FotoCheckKameraSamping   string `json:"fotoCheckKameraSamping"`
	FotoKunciCrm             string `json:"fotoKunciCrm"`
	FotoStrukErrorLog        string `json:"fotoStrukErrorLog"`
	FotoContactlessReader    string `json:"fotoContactlessReader"`
	FotoKomponenPc           string `json:"fotoKomponenPc"`
	FotoStikerBriit          string `json:"fotoStikerBriit"`
	TanggalPrestaging        string `json:"tanggalPrestaging"`
	TextNotes                string `json:"textNotes"`
}
