package salesResponseDTO

import "time"

type MISSuperAdmin struct {
	SalesName     string    `json:"salesName" gorm:"column:sales_name"`
	Metadata      string    `json:"metadata"`
	Status        string    `json:"status"`
	JenisProperti string    `json:"jenisProperti" gorm:"column:jenis_properti"`
	TipeProperti  string    `json:"tipeProperti" gorm:"column:tipe_properti"`
	CreatedAt     time.Time `json:"-"`
	ModifiedAt    time.Time `json:"-"`
	CreatedAtRes  string    `json:"createdAt" gorm:"-"`
	ModifiedAtRes string    `json:"modifiedAt" gorm:"-"`
}

type MISDeveloper struct {
	ID             int       `gorm:"primary_key:auto_increment" json:"-"`
	IDResponse     string    `json:"id" gorm:"-"`
	EmailDeveloper string    `gorm:"column:developer_email;size:255" json:"developerEmail"`
	EmailSales     string    `gorm:"column:sales_email;size:255" json:"salesEmail"`
	RefferalCode   string    `gorm:"column:refferal_code;size:255" json:"refferalCode"`
	RegisteredBy   string    `gorm:"column:registered_by;size:255" json:"registeredBy"`
	SalesName      string    `gorm:"column:sales_name;size:255" json:"salesName"`
	SalesPhone     string    `json:"salesPhone" gorm:"column:salesPhone"`
	Status         string    `json:"status" gorm:"column:status"`
	CreatedAt      time.Time `json:"-"`
	ModifiedAt     time.Time `json:"-"`
	CreatedAtRes   string    `json:"createdAt" gorm:"-"`
	ModifiedAtRes  string    `json:"modifiedAt" gorm:"-"`
}
