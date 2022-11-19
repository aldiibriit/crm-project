package entity

type Properti struct {
	ID                int    `json:"id"`
	GroupProperti     string `json:"groupProperti" gorm:"column:group_properti"`
	Email             string `json:"email"`
	NamaProperti      string `json:"namaProperti" gorm:"column:nama_properti"`
	DeskripsiProperti string `json:"deskripsiProperti" gorm:"column:deskripsi_properti"`
	HargaProperti     string `json:"hargaProperti" gorm:"column:harga_properti"`
	JumlahLantai      string `json:"jmlLantai" gorm:"column:jml_lantai"`
	ProjectId         int    `json:"projectId" gorm:"column:project_id"`
	ClusterId         int    `json:"clusterId" gorm:"column:cluster_id"`
}
