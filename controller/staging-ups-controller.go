package controller

import (
	"errors"
	"go-api/helper"
	"go-api/service"
	"mime/multipart"
	"net/http"
	"os"

	stagingUpsRequest "go-api/dto/request/staging-ups"
	"go-api/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StagingUPSController interface {
	PostStagingUPS(ctx *gin.Context)
	ApproveStagingUPS(ctx *gin.Context)
	RejectStagingUPS(ctx *gin.Context)
	ReuploadStagingUPS(ctx *gin.Context)
	AllSubmittedDataStagingUPS(ctx *gin.Context)
	GetSubmittedDataStagingUPSBySn(ctx *gin.Context)
	GetRejectedDataStagingUPS(ctx *gin.Context)
}

type stagingUPSController struct {
	stagingUpsService service.StagingUPSService
	jwtService        service.JWTService
}

func NewStagingUPSController(stagingUpsServ service.StagingUPSService, jwtServ service.JWTService) StagingUPSController {
	return &stagingUPSController{
		stagingUpsService: stagingUpsServ,
		jwtService:        jwtServ,
	}
}

func (controller *stagingUPSController) PostStagingUPS(ctx *gin.Context) {
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
	fotoUpsFull, _ := ctx.FormFile("fotoUpsFull")
	fotoKelengkapanUps, _ := ctx.FormFile("fotoKelengkapanUps")
	fotoStikerBitDanSucofindo, _ := ctx.FormFile("fotoStikerBitDanSucofindo")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoUpsFull.Filename] = fotoUpsFull
	requestMap[fotoKelengkapanUps.Filename] = fotoKelengkapanUps
	requestMap[fotoStikerBitDanSucofindo.Filename] = fotoStikerBitDanSucofindo

	request := stagingUpsRequest.PostStaging{
		IdUploader:                idUploader,
		Uploader:                  uploader,
		Sn:                        sn,
		ProjectName:               projectName,
		FotoUpsFull:               fotoUpsFull.Filename,
		FotoKelengkapanUps:        fotoKelengkapanUps.Filename,
		FotoStikerBitDanSucofindo: fotoStikerBitDanSucofindo.Filename,
		TextNotes:                 textNotes,
	}

	response = controller.stagingUpsService.PostStaging(requestMap, request)
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

func (controller *stagingUPSController) ApproveStagingUPS(ctx *gin.Context) {
	var request stagingUpsRequest.ApproveStaging
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
	response = controller.stagingUpsService.ApproveStaging(request)
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

func (controller *stagingUPSController) RejectStagingUPS(ctx *gin.Context) {
	var request stagingUpsRequest.RejectStaging
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
	response = controller.stagingUpsService.RejectStaging(request)
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

func (controller *stagingUPSController) ReuploadStagingUPS(ctx *gin.Context) {
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
	fotoUpsFull, _ := ctx.FormFile("fotoUpsFull")
	fotoKelengkapanUps, _ := ctx.FormFile("fotoKelengkapanUps")
	fotoStikerBitDanSucofindo, _ := ctx.FormFile("fotoStikerBitDanSucofindo")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoUpsFull.Filename] = fotoUpsFull
	requestMap[fotoKelengkapanUps.Filename] = fotoKelengkapanUps
	requestMap[fotoStikerBitDanSucofindo.Filename] = fotoStikerBitDanSucofindo

	request := stagingUpsRequest.PostStaging{
		IdUploader:                idUploader,
		Uploader:                  uploader,
		Sn:                        sn,
		ProjectName:               projectName,
		FotoUpsFull:               fotoUpsFull.Filename,
		FotoKelengkapanUps:        fotoKelengkapanUps.Filename,
		FotoStikerBitDanSucofindo: fotoStikerBitDanSucofindo.Filename,
		TextNotes:                 textNotes,
	}

	response = controller.stagingUpsService.ReuploadStaging(requestMap, request)
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

func (controller *stagingUPSController) AllSubmittedDataStagingUPS(ctx *gin.Context) {
	response := controller.stagingUpsService.AllSubmittedData()
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingUPSController) GetSubmittedDataStagingUPSBySn(ctx *gin.Context) {
	var request stagingUpsRequest.FindBySn
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
	response = controller.stagingUpsService.GetSubmittedDataBySn(request)
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingUPSController) GetRejectedDataStagingUPS(ctx *gin.Context) {
	var request stagingUpsRequest.FindRejectedData
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
	response = controller.stagingUpsService.GetRejectedData(request)
	ctx.JSON(response.HttpCode, response)
}
