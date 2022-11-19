package entity

import "time"

type TblSales struct {
	ID             int    `gorm:"primary_key:auto_increment" json:"id"`
	EmailDeveloper string `gorm:"column:developer_email;size:255" json:"developerEmail"`
	EmailSales     string `gorm:"column:sales_email;size:255" json:"salesEmail"`
	RefferalCode   string `gorm:"column:refferal_code;size:255" json:"refferalCode"`
	RegisteredBy   string `gorm:"column:registered_by;size:255" json:"registeredBy"`
	CreatedAt      time.Time
	ModifiedAt     time.Time
	Customers      []TblCustomer `gorm:"foreignKey:SalesID;references:ID"`
}
