package userRequestDTO

type ListUserReferralRequestDTO struct {
	SalesEmail string `json:"salesEmail"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}
