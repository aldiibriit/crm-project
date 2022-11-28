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
	ListProject(ctx *gin.Context)
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
	var request salesRequestDTO.MISDeveloperRequestDTO

	errDTO := ctx.ShouldBind(&request)
	if errDTO != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		response.ResponseData = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	decryptedRequest, err := deserializeMisDeveloperRequest(request)
	if err != nil {
		response.HttpCode = 400
		response.ResponseCode = "99"
		response.ResponseDesc = "Error in deserialize"
		response.ResponseData = nil
		response.Summary = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	response = controller.salesService.MISDeveloper(decryptedRequest)
	ctx.JSON(response.HttpCode, response)
}

func (controller *salesController) MISSuperAdmin(ctx *gin.Context) {
	var response responseDTO.Response
	var request salesRequestDTO.MISSuperAdminRequestDTO
	errDTO := ctx.ShouldBind(&request)
	if errDTO != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		response.ResponseData = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	response = controller.salesService.MISSuperAdmin(request)
	ctx.JSON(response.HttpCode, response)
}

func (controller *salesController) ListProject(ctx *gin.Context) {
	var response responseDTO.Response
	var request salesRequestDTO.ListProjectRequest
	errDTO := ctx.ShouldBind(&request)
	if errDTO != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		response.ResponseData = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	decryptedRequest, err := deserializeListProjectBySales(request)
	if err != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "Deserialize error ! " + errDTO.Error()
		response.Summary = nil
		response.ResponseData = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	response = controller.salesService.ListProject(decryptedRequest)
	ctx.JSON(response.HttpCode, response)
}

func deserializeMisDeveloperRequest(request interface{}) (salesRequestDTO.MISDeveloperRequestDTO, error) {
	otpDTO := request.(salesRequestDTO.MISDeveloperRequestDTO)

	cipheTextEmailDeveloper, _ := base64.StdEncoding.DecodeString(otpDTO.EmailDeveloper)
	plainTextEmailDeveloper, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmailDeveloper))
	if err != nil {
		return salesRequestDTO.MISDeveloperRequestDTO{}, err
	}

	var result salesRequestDTO.MISDeveloperRequestDTO

	result.EmailDeveloper = plainTextEmailDeveloper
	result.Keyword = otpDTO.Keyword
	result.Offset = otpDTO.Offset
	result.Limit = otpDTO.Limit

	return result, nil
}

func deserializeListProjectBySales(request interface{}) (salesRequestDTO.ListProjectRequest, error) {
	otpDTO := request.(salesRequestDTO.ListProjectRequest)

	cipheTextEmailSales, _ := base64.StdEncoding.DecodeString(otpDTO.EmailSales)
	plainTextEmailSales, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmailSales))
	if err != nil {
		return salesRequestDTO.ListProjectRequest{}, err
	}

	var result salesRequestDTO.ListProjectRequest

	result.EmailSales = plainTextEmailSales
	result.PageStart = otpDTO.PageStart

	return result, nil
}
