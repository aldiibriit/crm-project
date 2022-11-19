package authRequestDTO

type ActivateRequestDTO struct {
	Email          string `json:"email"`
	UrlEncoded     string `json:"urlEncoded"`
	RegistrationId string `json:"registrationID"`
	Action         int    `json:"action"`
}
