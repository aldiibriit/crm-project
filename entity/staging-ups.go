package entity

type TbStagingUps struct {
	ID          uint64 `json:"id" gorm:"primary_key:auto_increment" `
	IdUploader  string `json:"idUploader"`
	Uploader    string `json:"uploader"`
	Sn          string `json:"sn"`
	ProjectName string `json:"projectName"`
}
