package salesRequestDTO

type ListProjectRequest struct {
	EmailSales   string `json:"emailSales"`
	ReferralCode string `json:"referralCode"`
	PageStart    int    `json:"pageStart"`
}
