package entity

type TbLogActivity struct {
	Category          string `json:"category" gorm:"column:category"`
	Sn                string `json:"sm" gorm:"column:sn"`
	StatusDescription string `json:"statusDescription" gorm:"column:status_description"`
}
