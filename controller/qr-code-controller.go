package controller

import (
	"errors"
	QrCodeRequest "go-api/dto/request/qr-code"
	"go-api/dto/response"
	"go-api/helper"
	"go-api/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type QrCodeController interface {
	GenerateQr(ctx *gin.Context)
}

type qrCodeController struct {
	qrCodeService service.QrCodeService
}

func NewQrCodeController(qrCodeServ service.QrCodeService) QrCodeController {
	return &qrCodeController{
		qrCodeService: qrCodeServ,
	}
}

func (controller *qrCodeController) GenerateQr(ctx *gin.Context) {
	var request QrCodeRequest.GenerateQr
	var badRequestResponse response.BadRequestResponse
	var response response.UniversalResponse
	// var badRequestResponse response.BadRequestResponse
	defer func() {
		if r := recover(); r != nil {
			response.HttpCode = 500
			response.ResponseCode = "99"
			response.ResponseMessage = os.Getenv("UNCHAUGHT_ERROR_MESSAGE")
			response.Data = nil
			ctx.JSON(response.HttpCode, response)
			return
		}
	}()

	err := ctx.ShouldBind(&request)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := helper.CustomValidator(ve)
			badRequestResponse.HttpCode = 400
			badRequestResponse.ResponseCode = "99"
			badRequestResponse.ResponseMessage = os.Getenv("INVALID_REQUEST_MESSAGE")
			badRequestResponse.Errors = out
			ctx.AbortWithStatusJSON(http.StatusBadRequest, badRequestResponse)
			return
		}
	}
	response = controller.qrCodeService.GenerateQr(request)
	ctx.JSON(response.HttpCode, response)
}
