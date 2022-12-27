package KPRRequestDTO

type ListPengajuanKPR struct {
	Email string `json:"email" binding:"required"`
}
