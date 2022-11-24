package salesRequestDTO

type MISSuperAdminRequestDTO struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit" binding:"required"`
	Offset  int    `json:"offset" binding:"required"`
}
