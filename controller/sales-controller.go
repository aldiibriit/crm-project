package controller

import (
	"encoding/base64"
	"fmt"
	"go-api/dto/request/salesRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/helper"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

type SalesController interface {
	MISDeveloper(ctx *gin.Context)
	MISSuperAdmin(ctx *gin.Context)
	DetailSalesByDeveloper(ctx *gin.Context)
	EditSalesByDeveloper(ctx *gin.Context)
	DeleteSalesByDeveloper(ctx *gin.Context)
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
		response.ResponseDesc = "Deserialize error ! "
		response.Summary = nil
		response.ResponseData = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	response = controller.salesService.ListProject(decryptedRequest)
	ctx.JSON(response.HttpCode, response)
}

func (controller *salesController) EditSalesByDeveloper(ctx *gin.Context) {
	var response responseDTO.Response
	var request salesRequestDTO.SalesEditRequestDTO

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

	decryptedRequest, err := deserializeEditSalesByDeveloper(request)
	if err != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		response.ResponseData = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}
	response = controller.salesService.EditSalesByDeveloper(decryptedRequest)

	ctx.JSON(response.HttpCode, response)
}

func (controller *salesController) DeleteSalesByDeveloper(ctx *gin.Context) {
	var response responseDTO.Response
	var request salesRequestDTO.SalesDeleteRequestDTO

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

	decryptedRequest, err := deserializeDeleteSales(request)
	if err != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		response.ResponseData = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	response = controller.salesService.DeleteSalesByDeveloper(decryptedRequest)

	ctx.JSON(response.HttpCode, response)
}

func (controller *salesController) DetailSalesByDeveloper(ctx *gin.Context) {
	var response responseDTO.Response
	var request salesRequestDTO.DetailSalesRequest

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

	decryptedRequest, err := deserializeDetailSales(request)
	if err != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		response.ResponseData = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	response = controller.salesService.DetailSalesByDeveloper(decryptedRequest)

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

	cipheTextEmailSales, err := base64.StdEncoding.DecodeString(otpDTO.EmailSales)
	if err != nil {
		fmt.Println(err.Error())
	}
	plainTextEmailSales, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmailSales))
	if err != nil {
		fmt.Println(err.Error())
		return salesRequestDTO.ListProjectRequest{}, err
	}

	var result salesRequestDTO.ListProjectRequest

	result.EmailSales = plainTextEmailSales
	result.PageStart = otpDTO.PageStart

	return result, nil
}

func deserializeEditSalesByDeveloper(request interface{}) (salesRequestDTO.SalesEditRequestDTO, error) {
	otpDTO := request.(salesRequestDTO.SalesEditRequestDTO)

	cipheTextEmail, err := base64.StdEncoding.DecodeString(otpDTO.Email)
	if err != nil {
		fmt.Println(err.Error())
	}
	cipheTextID, err := base64.StdEncoding.DecodeString(otpDTO.ID)
	if err != nil {
		fmt.Println(err.Error())
	}
	cipheTextSalesName, err := base64.StdEncoding.DecodeString(otpDTO.SalesName)
	if err != nil {
		fmt.Println(err.Error())
	}
	cipheTextSalesPhone, err := base64.StdEncoding.DecodeString(otpDTO.SalesPhone)
	if err != nil {
		fmt.Println(err.Error())
	}

	plainTextEmail, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmail))
	if err != nil {
		fmt.Println(err.Error())
		return salesRequestDTO.SalesEditRequestDTO{}, err
	}

	plainTextID, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextID))
	if err != nil {
		fmt.Println(err.Error())
		return salesRequestDTO.SalesEditRequestDTO{}, err
	}

	plainTextSalesName, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextSalesName))
	if err != nil {
		fmt.Println(err.Error())
		return salesRequestDTO.SalesEditRequestDTO{}, err
	}

	plainTextSalesPhone, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextSalesPhone))
	if err != nil {
		fmt.Println(err.Error())
		return salesRequestDTO.SalesEditRequestDTO{}, err
	}

	var result salesRequestDTO.SalesEditRequestDTO

	result.Email = plainTextEmail
	result.ID = plainTextID
	result.SalesName = plainTextSalesName
	result.SalesPhone = plainTextSalesPhone

	return result, nil
}

func deserializeDetailSales(request interface{}) (salesRequestDTO.DetailSalesRequest, error) {
	otpDTO := request.(salesRequestDTO.DetailSalesRequest)

	cipheTextID, err := base64.StdEncoding.DecodeString(otpDTO.ID)
	if err != nil {
		fmt.Println(err.Error())
	}
	plainTextID, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextID))
	if err != nil {
		fmt.Println(err.Error())
		return salesRequestDTO.DetailSalesRequest{}, err
	}

	var result salesRequestDTO.DetailSalesRequest

	result.ID = plainTextID

	return result, nil
}

func deserializeDeleteSales(request interface{}) (salesRequestDTO.SalesDeleteRequestDTO, error) {
	otpDTO := request.(salesRequestDTO.SalesDeleteRequestDTO)

	cipheTextID, err := base64.StdEncoding.DecodeString(otpDTO.ID)
	if err != nil {
		fmt.Println(err.Error())
	}
	plainTextID, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextID))
	if err != nil {
		fmt.Println(err.Error())
		return salesRequestDTO.SalesDeleteRequestDTO{}, err
	}

	var result salesRequestDTO.SalesDeleteRequestDTO

	result.ID = plainTextID

	return result, nil
}
