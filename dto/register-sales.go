package dto

type RegisterSalesDTO struct {
	EmailDeveloper string `json:"emailDeveloper" binding:"required"`
	EmailSales     string `json:"emailSales" binding:"required"`
	SalesName      string `json:"salesName" binding:"required"`
	SalesPhone     string `json:"salesPhone" binding:"required"`
	RegisteredBy   string `json:"registeredBy" binding:"required"`
}
