package controller

import (
	"errors"
	"go-api/helper"
	"go-api/service"
	"mime/multipart"
	"net/http"
	"os"

	stagingDigitalSignageRequest "go-api/dto/request/staging-digital-signage"
	"go-api/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StagingDigitalSignageController interface {
	PostStagingDigitalSignage(ctx *gin.Context)
	ApproveStagingDigitalSignage(ctx *gin.Context)
	RejectStagingDigitalSignage(ctx *gin.Context)
	ReuploadStagingDigitalSignage(ctx *gin.Context)
	AllSubmittedDataStagingDigitalSignage(ctx *gin.Context)
	GetSubmittedDataStagingDigitalSignageBySn(ctx *gin.Context)
	GetRejectedDataStagingDigitalSignage(ctx *gin.Context)
}

type stagingDigitalSignageController struct {
	stagingDigitalSignageService service.StagingDigitalSignageService
	jwtService                   service.JWTService
}

func NewStagingDigitalSignageController(stagingDigitalSignageServ service.StagingDigitalSignageService, jwtServ service.JWTService) StagingDigitalSignageController {
	return &stagingDigitalSignageController{
		stagingDigitalSignageService: stagingDigitalSignageServ,
		jwtService:                   jwtServ,
	}
}

func (controller *stagingDigitalSignageController) PostStagingDigitalSignage(ctx *gin.Context) {
	var response response.UniversalResponse
	defer func() {
		if r := recover(); r != nil {
			response.HttpCode = 500
			response.ResponseCode = "99"
			response.ResponseMessage = "Unchaught error"
			response.Data = nil
			ctx.AbortWithStatusJSON(response.HttpCode, response)
			return
		}
	}()

	idUploader := ctx.PostForm("idUploader")
	uploader := ctx.PostForm("uploader")
	sn := ctx.PostForm("sn")
	projectName := ctx.PostForm("projectName")
	textNotes := ctx.PostForm("textNotes")
	fotoLed, _ := ctx.FormFile("fotoLed")
	fotoSnLedDus, _ := ctx.FormFile("fotoSnLedDus")
	fotoStikerBitDanSucofindo, _ := ctx.FormFile("fotoStikerBitDanSucofindo")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoLed.Filename] = fotoLed
	requestMap[fotoSnLedDus.Filename] = fotoSnLedDus
	requestMap[fotoStikerBitDanSucofindo.Filename] = fotoStikerBitDanSucofindo

	request := stagingDigitalSignageRequest.PostStaging{
		IdUploader:                idUploader,
		Uploader:                  uploader,
		Sn:                        sn,
		ProjectName:               projectName,
		FotoLed:                   fotoLed.Filename,
		FotoSnLedDus:              fotoSnLedDus.Filename,
		FotoStikerBitDanSucofindo: fotoStikerBitDanSucofindo.Filename,
		TextNotes:                 textNotes,
	}

	response = controller.stagingDigitalSignageService.PostStaging(requestMap, request)
	if response.ResponseCode == "" || response.ResponseMessage == "" {
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("UNCHAUGHT_ERROR_MESSAGE")
		response.Data = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingDigitalSignageController) ApproveStagingDigitalSignage(ctx *gin.Context) {
	var request stagingDigitalSignageRequest.ApproveStaging
	var badRequestResponse response.BadRequestResponse
	var response response.UniversalResponse

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
	response = controller.stagingDigitalSignageService.ApproveStaging(request)
	if response.ResponseCode == "" || response.ResponseMessage == "" {
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("UNCHAUGHT_ERROR_MESSAGE")
		response.Data = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingDigitalSignageController) RejectStagingDigitalSignage(ctx *gin.Context) {
	var request stagingDigitalSignageRequest.RejectStaging
	var badRequestResponse response.BadRequestResponse
	var response response.UniversalResponse

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
	response = controller.stagingDigitalSignageService.RejectStaging(request)
	if response.ResponseCode == "" || response.ResponseMessage == "" {
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("UNCHAUGHT_ERROR_MESSAGE")
		response.Data = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingDigitalSignageController) ReuploadStagingDigitalSignage(ctx *gin.Context) {
	var response response.UniversalResponse
	defer func() {
		if r := recover(); r != nil {
			response.HttpCode = 500
			response.ResponseCode = "99"
			response.ResponseMessage = "Unchaught error"
			response.Data = nil
			ctx.AbortWithStatusJSON(response.HttpCode, response)
			return
		}
	}()

	idUploader := ctx.PostForm("idUploader")
	uploader := ctx.PostForm("uploader")
	sn := ctx.PostForm("sn")
	projectName := ctx.PostForm("projectName")
	textNotes := ctx.PostForm("textNotes")
	fotoLed, _ := ctx.FormFile("fotoLed")
	fotoSnLedDus, _ := ctx.FormFile("fotoSnLedDus")
	fotoStikerBitDanSucofindo, _ := ctx.FormFile("fotoStikerBitDanSucofindo")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoLed.Filename] = fotoLed
	requestMap[fotoSnLedDus.Filename] = fotoSnLedDus
	requestMap[fotoStikerBitDanSucofindo.Filename] = fotoStikerBitDanSucofindo

	request := stagingDigitalSignageRequest.PostStaging{
		IdUploader:                idUploader,
		Uploader:                  uploader,
		Sn:                        sn,
		ProjectName:               projectName,
		FotoLed:                   fotoLed.Filename,
		FotoSnLedDus:              fotoSnLedDus.Filename,
		FotoStikerBitDanSucofindo: fotoStikerBitDanSucofindo.Filename,
		TextNotes:                 textNotes,
	}

	response = controller.stagingDigitalSignageService.ReuploadStaging(requestMap, request)
	if response.ResponseCode == "" || response.ResponseMessage == "" {
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("UNCHAUGHT_ERROR_MESSAGE")
		response.Data = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingDigitalSignageController) AllSubmittedDataStagingDigitalSignage(ctx *gin.Context) {
	response := controller.stagingDigitalSignageService.AllSubmittedData()
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingDigitalSignageController) GetSubmittedDataStagingDigitalSignageBySn(ctx *gin.Context) {
	var request stagingDigitalSignageRequest.FindBySn
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
	response = controller.stagingDigitalSignageService.GetSubmittedDataBySn(request)
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingDigitalSignageController) GetRejectedDataStagingDigitalSignage(ctx *gin.Context) {
	var request stagingDigitalSignageRequest.FindRejectedData
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
	response = controller.stagingDigitalSignageService.GetRejectedData(request)
	ctx.JSON(response.HttpCode, response)
}
