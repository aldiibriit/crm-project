package controller

import (
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

	res := controller.otpService.ValidateOTP(validateOTPRequest)
	ctx.JSON(res.HttpCode, res)
}
