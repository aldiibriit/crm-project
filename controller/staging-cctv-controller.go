package controller

import (
	"errors"
	"go-api/helper"
	"go-api/service"
	"mime/multipart"
	"net/http"
	"os"

	stagingCctvRequest "go-api/dto/request/staging-cctv"
	"go-api/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StagingCCTVController interface {
	PostStagingCCTV(ctx *gin.Context)
	ApproveStagingCCTV(ctx *gin.Context)
	RejectStagingCCTV(ctx *gin.Context)
	ReuploadStagingCCTV(ctx *gin.Context)
	AllSubmittedDataStagingCCTV(ctx *gin.Context)
	GetSubmittedDataStagingCCTVBySn(ctx *gin.Context)
	GetRejectedDataStagingCCTV(ctx *gin.Context)
}

type stagingCCTVController struct {
	stagingCctvService service.StagingCCTVService
	jwtService         service.JWTService
}

func NewStagingCCTVController(stagingCctvServ service.StagingCCTVService, jwtServ service.JWTService) StagingCCTVController {
	return &stagingCCTVController{
		stagingCctvService: stagingCctvServ,
		jwtService:         jwtServ,
	}
}

func (controller *stagingCCTVController) PostStagingCCTV(ctx *gin.Context) {
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
	fotoNvrDanCamera, _ := ctx.FormFile("fotoNvrDanCamera")
	fotoSnNvrDanCameraDus, _ := ctx.FormFile("fotoSnNvrDanCameraDus")
	fotoStikerBitDanSucofindo, _ := ctx.FormFile("fotoStikerBitDanSucofindo")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoNvrDanCamera.Filename] = fotoNvrDanCamera
	requestMap[fotoSnNvrDanCameraDus.Filename] = fotoSnNvrDanCameraDus
	requestMap[fotoStikerBitDanSucofindo.Filename] = fotoStikerBitDanSucofindo

	request := stagingCctvRequest.PostStaging{
		IdUploader:                idUploader,
		Uploader:                  uploader,
		Sn:                        sn,
		ProjectName:               projectName,
		FotoNvrDanCamera:          fotoNvrDanCamera.Filename,
		FotoSnNvrDanCameraDus:     fotoSnNvrDanCameraDus.Filename,
		FotoStikerBitDanSucofindo: fotoStikerBitDanSucofindo.Filename,
		TextNotes:                 textNotes,
	}

	response = controller.stagingCctvService.PostStaging(requestMap, request)
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

func (controller *stagingCCTVController) ApproveStagingCCTV(ctx *gin.Context) {
	var request stagingCctvRequest.ApproveStaging
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
	response = controller.stagingCctvService.ApproveStaging(request)
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

func (controller *stagingCCTVController) RejectStagingCCTV(ctx *gin.Context) {
	var request stagingCctvRequest.RejectStaging
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
	response = controller.stagingCctvService.RejectStaging(request)
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

func (controller *stagingCCTVController) ReuploadStagingCCTV(ctx *gin.Context) {
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
	fotoNvrDanCamera, _ := ctx.FormFile("fotoNvrDanCamera")
	fotoSnNvrDanCameraDus, _ := ctx.FormFile("fotoSnNvrDanCameraDus")
	fotoStikerBitDanSucofindo, _ := ctx.FormFile("fotoStikerBitDanSucofindo")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoNvrDanCamera.Filename] = fotoNvrDanCamera
	requestMap[fotoSnNvrDanCameraDus.Filename] = fotoSnNvrDanCameraDus
	requestMap[fotoStikerBitDanSucofindo.Filename] = fotoStikerBitDanSucofindo

	request := stagingCctvRequest.PostStaging{
		IdUploader:                idUploader,
		Uploader:                  uploader,
		Sn:                        sn,
		ProjectName:               projectName,
		FotoNvrDanCamera:          fotoNvrDanCamera.Filename,
		FotoSnNvrDanCameraDus:     fotoSnNvrDanCameraDus.Filename,
		FotoStikerBitDanSucofindo: fotoStikerBitDanSucofindo.Filename,
		TextNotes:                 textNotes,
	}

	response = controller.stagingCctvService.ReuploadStaging(requestMap, request)
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

func (controller *stagingCCTVController) AllSubmittedDataStagingCCTV(ctx *gin.Context) {
	response := controller.stagingCctvService.AllSubmittedData()
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingCCTVController) GetSubmittedDataStagingCCTVBySn(ctx *gin.Context) {
	var request stagingCctvRequest.FindBySn
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
	response = controller.stagingCctvService.GetSubmittedDataBySn(request)
	ctx.JSON(response.HttpCode, response)
}

func (controller *stagingCCTVController) GetRejectedDataStagingCCTV(ctx *gin.Context) {
	var request stagingCctvRequest.FindRejectedData
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
	response = controller.stagingCctvService.GetRejectedData(request)
	ctx.JSON(response.HttpCode, response)
}
