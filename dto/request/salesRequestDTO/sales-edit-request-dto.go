package salesRequestDTO

type SalesEditRequestDTO struct {
	ID         string `json:"id" binding:"required"`
	Email      string `json:"email" binding:"required"`
	SalesName  string `json:"salesName" binding:"required"`
	SalesPhone string `json:"salesPhone" binding:"required"`
}
