package salesRequestDTO

type MISSuperAdminRequestDTO struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit" binding:"required"`
	Offset  int    `json:"offset" binding:"required"`
}

type MISDeveloperRequestDTO struct {
	EmailDeveloper string `json:"emailDeveloper" binding:"required"`
	Keyword        string `json:"keyword"`
	Limit          string `json:"limit" binding:"required"`
	Offset         string `json:"offset" binding:"required"`
}
