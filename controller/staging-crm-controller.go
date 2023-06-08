package controller

import (
	"errors"
	"go-api/helper"
	"go-api/service"
	"mime/multipart"
	"net/http"
	"os"

	stagingCrmRequest "go-api/dto/request/staging-crm"
	"go-api/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StagingCRMController interface {
	PostStagingCRM(ctx *gin.Context)
	ApproveStagingCRM(ctx *gin.Context)
	RejectStagingCRM(ctx *gin.Context)
	ReuploadStagingCRM(ctx *gin.Context)
	AllSubmittedDataStagingCRM(ctx *gin.Context)
	GetSubmittedDataStagingCRMBySn(ctx *gin.Context)
	GetRejectedDataStagingCRM(ctx *gin.Context)
}

type stagingCRMController struct {
	stagingCrmService service.StagingCRMService
	jwtService        service.JWTService
}

func NewStagingCRMController(stagingCrmServ service.StagingCRMService, jwtServ service.JWTService) StagingCRMController {
	return &stagingCRMController{
		stagingCrmService: stagingCrmServ,
		jwtService:        jwtServ,
	}
}

func (controller *stagingCRMController) PostStagingCRM(ctx *gin.Context) {
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
	fotoStikerBitDanSucofindo, _ := ctx.FormFile("fotoStikerBitDanSucofindo")
	fotoMesinCrmFull, _ := ctx.FormFile("fotoMesinCrmFull")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoStikerBitDanSucofindo.Filename] = fotoStikerBitDanSucofindo
	requestMap[fotoMesinCrmFull.Filename] = fotoMesinCrmFull

	request := stagingCrmRequest.PostStaging{
		IdUploader:                idUploader,
		Uploader:                  uploader,
		Sn:                        sn,
		ProjectName:               projectName,
		FotoMesinCrmFull:          fotoMesinCrmFull.Filename,
		FotoStikerBitDanSucofindo: fotoStikerBitDanSucofindo.Filename,
		TextNotes:                 textNotes,
	}

	response = controller.stagingCrmService.PostStaging(requestMap, request)
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

func (controller *stagingCRMController) ApproveStagingCRM(ctx *gin.Context) {
	var request stagingCrmRequest.ApproveStaging
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
	response = controller.stagingCrmService.ApproveStaging(request)
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

func (controller *stagingCRMController) RejectStagingCRM(ctx *gin.Context) {
	var request stagingCrmRequest.RejectStaging
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
	response = controller.stagingCrmService.RejectStaging(request)
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

func (controller *stagingCRMController) ReuploadStagingCRM(ctx *gin.Context) {
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
	fotoMesinCrmFull, _ := ctx.FormFile("fotoMesinCrmFull")
	fotoStikerBitDanSucofindo, _ := ctx.FormFile("fotoStikerBitDanSucofindo")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoMesinCrmFull.Filename] = fotoMesinCrmFull
	requestMap[fotoStikerBitDanSucofindo.Filename] = fotoStikerBitDanSucofindo

	request := stagingCrmRequest.PostStaging{
		IdUploader:                idUploader,
		Uploader:                  uploader,
		Sn:                        sn,
		ProjectName:               projectName,
		FotoMesinCrmFull:          fotoMesinCrmFull.Filename,
		FotoStikerBitDanSucofindo: fotoStikerBitDanSucofindo.Filename,
		TextNotes:                 textNotes,
	}

	response = controller.stagingCrmService.ReuploadStaging(requestMap, request)
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

func (controller *stagingCRMController) AllSubmittedDataStagingCRM(ctx *gin.Context) {
	response := controller.stagingCrmService.AllSubmittedData()
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingCRMController) GetSubmittedDataStagingCRMBySn(ctx *gin.Context) {
	var request stagingCrmRequest.FindBySn
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
	response = controller.stagingCrmService.GetSubmittedDataBySn(request)
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingCRMController) GetRejectedDataStagingCRM(ctx *gin.Context) {
	var request stagingCrmRequest.FindRejectedData
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
	response = controller.stagingCrmService.GetRejectedData(request)
	ctx.JSON(response.HttpCode, response)
}
