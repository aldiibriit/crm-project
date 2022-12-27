package controller

import (
	"go-api/dto/request/KPRRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

type KPRController interface {
	PengajuanKPR(ctx *gin.Context)
	ListPengajuan(ctx *gin.Context)
}

type kprController struct {
	kprService service.KPRService
}

func NewKPRController(kprServ service.KPRService) KPRController {
	return &kprController{
		kprService: kprServ,
	}
}

func (c *kprController) PengajuanKPR(ctx *gin.Context) {
	var response responseDTO.Response
	var request KPRRequestDTO.PengajuanKPRRequest
	errDTO := ctx.ShouldBind(&request)
	if errDTO != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	response = c.kprService.PengajuanKPR(request)
	ctx.JSON(response.HttpCode, response)
}

func (c *kprController) ListPengajuan(ctx *gin.Context) {
	var request KPRRequestDTO.ListPengajuanKPR
	var response responseDTO.Response

	errDTO := ctx.ShouldBind(&request)
	if errDTO != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
	}

	response = c.kprService.ListPengajuan(request)
	ctx.JSON(response.HttpCode, response)
}
