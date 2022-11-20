package controller

import (
	"encoding/base64"
	"go-api/helper"
	"go-api/service"
	"net/http"

	"go-api/dto/request/otpRequestDTO"

	"github.com/gin-gonic/gin"
)

type OTPController interface {
	ValidateOTP(ctx *gin.Context)
}

type otpController struct {
	otpService service.OTPService
}

func NEwOTPController(otpServ service.OTPService) OTPController {
	return &otpController{
		otpService: otpServ,
	}
}

func (controller *otpController) ValidateOTP(ctx *gin.Context) {
	var validateOTPRequest otpRequestDTO.ValidateOTPRequest
	errDTO := ctx.ShouldBind(&validateOTPRequest)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	decryptedRequest := deserializeValidateOTPRequest(validateOTPRequest)

	res := controller.otpService.ValidateOTP(decryptedRequest)
	ctx.JSON(res.HttpCode, res)
}

func deserializeValidateOTPRequest(request interface{}) otpRequestDTO.ValidateOTPRequest {
	otpDTO := request.(otpRequestDTO.ValidateOTPRequest)

	cipheTextOTP, _ := base64.StdEncoding.DecodeString(otpDTO.OTP)
	cipheTextEmail, _ := base64.StdEncoding.DecodeString(otpDTO.Email)
	plainTextOTP, _ := helper.RsaDecryptFromFEInBE([]byte(cipheTextOTP))
	plainTextEmail, _ := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmail))

	var result otpRequestDTO.ValidateOTPRequest

	result.OTP = plainTextOTP
	result.Email = plainTextEmail

	return result
}
