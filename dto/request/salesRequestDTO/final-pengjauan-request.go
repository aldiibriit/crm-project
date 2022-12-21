package salesRequestDTO

type FinalPengajuanRequest struct {
	Email    string `json:"email" binding:"required"`
	UserType string `json:"userType" binding:"required"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}
