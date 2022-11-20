package dto

type PasswordConfirmationDTO struct {
	Email             string `json:"email"`
	NewPassword       string `json:"newPassword"`
	RetypeNewPassword string `json:"retypeNewPassword"`
}
