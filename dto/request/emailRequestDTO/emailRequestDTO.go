package emailRequestDTO

type EmailRequestDTO struct {
	Name       string `json:"name"`
	Subject    string `json:"subject"`
	Action     int    `json:"action"`
	EmailBody  string `json:"emailBody"`
	ToAddres   string `json:"toAddress"`
	UrlEncoded string `json:"urlEncoded"`
}
