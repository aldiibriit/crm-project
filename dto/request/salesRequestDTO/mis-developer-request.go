package salesRequestDTO

type AllRequest struct {
	EmailDeveloper string `json:"emailDeveloper" binding:"required"`
}
