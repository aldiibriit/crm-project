package salesRequestDTO

type ListProjectRequest struct {
	// EmailSales   string `json:"emailSales" binding:"required"`
	ReferralCode string `json:"referralCode" binding:"required"`
	PageStart    int    `json:"pageStart"`
}
