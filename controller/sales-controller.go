package controller

import (
	"encoding/base64"
	"go-api/dto/request/salesRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/helper"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

type SalesController interface {
	MISDeveloper(ctx *gin.Context)
}

type salesController struct {
	salesService service.SalesService
}

func NewSalesController(salesServ service.SalesService) SalesController {
	return &salesController{
		salesService: salesServ,
	}
}

func (controller *salesController) MISDeveloper(ctx *gin.Context) {
	var response responseDTO.Response
	var request salesRequestDTO.AllRequest
	errDTO := ctx.ShouldBind(&request)
	if errDTO != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		response.ResponseData = nil
		ctx.JSON(response.HttpCode, response)
	}
	decryptedRequest := deserializeAllMisDeveloperRequest(request)
	response = controller.salesService.MISDeveloper(decryptedRequest)
	ctx.JSON(response.HttpCode, response)
}

func deserializeAllMisDeveloperRequest(request interface{}) salesRequestDTO.AllRequest {
	otpDTO := request.(salesRequestDTO.AllRequest)

	cipheTextEmailDeveloper, _ := base64.StdEncoding.DecodeString(otpDTO.EmailDeveloper)
	plainTextEmailDeveloper, _ := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmailDeveloper))

	var result salesRequestDTO.AllRequest

	result.EmailDeveloper = plainTextEmailDeveloper

	return result
}
