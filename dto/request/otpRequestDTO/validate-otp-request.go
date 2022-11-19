package otpRequestDTO

type ValidateOTPRequest struct {
	OTP   string `json:"otp"`
	Email string `json:"email"`
}
