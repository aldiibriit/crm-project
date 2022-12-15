package userRequestDTO

type ListUserReferralRequestDTO struct {
	SalesEmail string `json:"salesEmail"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	Keyword    string `json:"keyword"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}
