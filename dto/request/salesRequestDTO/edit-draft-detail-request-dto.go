package salesRequestDTO

type EditDraftDetailRequestDTO struct {
	ID    string `json:"id" binding:"required"`
	NIK   string `json:"nik"`
	Email string `json:"email"`
}
