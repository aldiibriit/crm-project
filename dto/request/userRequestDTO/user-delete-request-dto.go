package userRequestDTO

type UserDeleteRequestDTO struct {
	ID string `json:"id" binding:"required"`
}
