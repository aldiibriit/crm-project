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
	MISSuperAdmin(ctx *gin.Context)
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
	decryptedRequest, err := deserializeAllMisDeveloperRequest(request)
	if err != nil {
		var response responseDTO.Response
		response.HttpCode = 400
		response.ResponseCode = "99"
		response.ResponseDesc = "Error in deserialize " + err.Error()
		response.ResponseData = nil
		response.Summary = nil
	}

	response = controller.salesService.MISDeveloper(decryptedRequest)
	ctx.JSON(response.HttpCode, response)
}

func (controller *salesController) MISSuperAdmin(ctx *gin.Context) {
	var response responseDTO.Response
	// var request salesRequestDTO.AllRequest
	// errDTO := ctx.ShouldBind(&request)
	// if errDTO != nil {
	// 	response.HttpCode = 400
	// 	response.MetadataResponse = nil
	// 	response.ResponseCode = "99"
	// 	response.ResponseDesc = errDTO.Error()
	// 	response.Summary = nil
	// 	response.ResponseData = nil
	// 	ctx.JSON(response.HttpCode, response)
	// }

	response = controller.salesService.MISSuperAdmin()
	ctx.JSON(response.HttpCode, response)
}

func deserializeAllMisDeveloperRequest(request interface{}) (salesRequestDTO.AllRequest, error) {
	otpDTO := request.(salesRequestDTO.AllRequest)

	cipheTextEmailDeveloper, _ := base64.StdEncoding.DecodeString(otpDTO.EmailDeveloper)
	plainTextEmailDeveloper, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmailDeveloper))
	if err != nil {
		return salesRequestDTO.AllRequest{}, err
	}

	var result salesRequestDTO.AllRequest

	result.EmailDeveloper = plainTextEmailDeveloper

	return result, nil
}
