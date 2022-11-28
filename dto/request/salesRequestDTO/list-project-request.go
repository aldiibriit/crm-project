package salesRequestDTO

type ListProjectRequest struct {
	EmailSales string `json:"emailSales" binding:"required"`
	PageStart  int    `json:"pageStart"`
}
