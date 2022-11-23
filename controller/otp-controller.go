package controller

import (
	"encoding/base64"
	"go-api/helper"
	"go-api/service"
	"net/http"

	"go-api/dto/request/otpRequestDTO"
	responseDTO "go-api/dto/response"

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

	decryptedRequest, err := deserializeValidateOTPRequest(validateOTPRequest)
	if err != nil {
		var response responseDTO.Response
		response.HttpCode = 400
		response.ResponseCode = "99"
		response.ResponseDesc = "Error in deserialize " + err.Error()
		response.ResponseData = nil
		response.Summary = nil
	}

	res := controller.otpService.ValidateOTP(decryptedRequest)
	ctx.JSON(res.HttpCode, res)
}

func deserializeValidateOTPRequest(request interface{}) (otpRequestDTO.ValidateOTPRequest, error) {
	otpDTO := request.(otpRequestDTO.ValidateOTPRequest)

	cipheTextOTP, _ := base64.StdEncoding.DecodeString(otpDTO.OTP)
	cipheTextEmail, _ := base64.StdEncoding.DecodeString(otpDTO.Email)
	plainTextOTP, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextOTP))
	if err != nil {
		return otpRequestDTO.ValidateOTPRequest{}, err
	}
	plainTextEmail, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmail))
	if err != nil {
		return otpRequestDTO.ValidateOTPRequest{}, err
	}

	var result otpRequestDTO.ValidateOTPRequest

	result.OTP = plainTextOTP
	result.Email = plainTextEmail

	return result, nil
}
