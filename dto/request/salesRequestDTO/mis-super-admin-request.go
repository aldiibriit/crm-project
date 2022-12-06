package salesRequestDTO

type MISSuperAdminRequestDTO struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit" `
	Offset  int    `json:"offset"`
}

type MISDeveloperRequestDTO struct {
	EmailDeveloper string `json:"emailDeveloper" binding:"required"`
	Keyword        string `json:"keyword"`
	Limit          int    `json:"limit" `
	Offset         int    `json:"offset"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
}
