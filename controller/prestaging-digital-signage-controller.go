package controller

import (
	"errors"
	"go-api/helper"
	"go-api/service"
	"mime/multipart"
	"net/http"
	"os"

	prestagingDigitalSignageRequest "go-api/dto/request/prestaging-digital-signage"
	"go-api/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PrestagingDigitalSignageController interface {
	PostPrestagingDigitalSignage(ctx *gin.Context)
	ApprovePrestagingDigitalSignage(ctx *gin.Context)
	RejectPrestagingDigitalSignage(ctx *gin.Context)
	ReuploadPrestagingDigitalSignage(ctx *gin.Context)
	AllSubmittedDataPrestagingDigitalSignage(ctx *gin.Context)
	GetSubmittedDataPrestagingDigitalSignageBySn(ctx *gin.Context)
	GetRejectedDataPrestagingDigitalSignage(ctx *gin.Context)
}

type prestagingDigitalSignageController struct {
	prestagingDigitalSignageService service.PrestagingDigitalSignageService
	jwtService                      service.JWTService
}

func NewPrestagingDigitalSignageController(prestagingDigitalSignageServe service.PrestagingDigitalSignageService, jwtServ service.JWTService) PrestagingDigitalSignageController {
	return &prestagingDigitalSignageController{
		prestagingDigitalSignageService: prestagingDigitalSignageServe,
		jwtService:                      jwtServ,
	}
}

func (controller *prestagingDigitalSignageController) PostPrestagingDigitalSignage(ctx *gin.Context) {
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
	snLed := ctx.PostForm("snLed")
	fotoLedFull, _ := ctx.FormFile("fotoLedFull")
	fotoSnLed, _ := ctx.FormFile("fotoSnLed")
	fotoTestDeadPixelPutih, _ := ctx.FormFile("fotoTestDeadPixelPutih")
	fotoTestDeadPixelBiru, _ := ctx.FormFile("fotoTestDeadPixelBiru")
	fotoTestDeadPixelMerah, _ := ctx.FormFile("fotoTestDeadPixelMerah")
	fotoStikerBit, _ := ctx.FormFile("fotoStikerBit")
	fotoTestDeadPixelHitam, _ := ctx.FormFile("fotoTestDeadPixelHitam")
	fotoKelengkapan, _ := ctx.FormFile("fotoKelengkapan")
	fotoTestDeadPixelHijau, _ := ctx.FormFile("fotoTestDeadPixelHijau")
	statusBarang := ctx.PostForm("statusBarang")
	brand := ctx.PostForm("brand")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoLedFull.Filename] = fotoLedFull
	requestMap[fotoSnLed.Filename] = fotoSnLed
	requestMap[fotoTestDeadPixelPutih.Filename] = fotoTestDeadPixelPutih
	requestMap[fotoTestDeadPixelBiru.Filename] = fotoTestDeadPixelBiru
	requestMap[fotoTestDeadPixelMerah.Filename] = fotoTestDeadPixelMerah
	requestMap[fotoStikerBit.Filename] = fotoStikerBit
	requestMap[fotoTestDeadPixelHitam.Filename] = fotoTestDeadPixelHitam
	requestMap[fotoKelengkapan.Filename] = fotoKelengkapan
	requestMap[fotoTestDeadPixelHijau.Filename] = fotoTestDeadPixelHijau

	request := prestagingDigitalSignageRequest.PostPrestaging{
		IdUploader:             idUploader,
		Uploader:               uploader,
		Sn:                     sn,
		ProjectName:            projectName,
		FotoLedFull:            fotoLedFull.Filename,
		FotoSnLed:              fotoSnLed.Filename,
		FotoTestDeadPixelPutih: fotoTestDeadPixelPutih.Filename,
		FotoTestDeadPixelBiru:  fotoTestDeadPixelBiru.Filename,
		FotoTestDeadPixelMerah: fotoTestDeadPixelMerah.Filename,
		FotoTestDeadPixelHijau: fotoTestDeadPixelHijau.Filename,
		FotoTestDeadPixelHitam: fotoTestDeadPixelHitam.Filename,
		FotoStikerBit:          fotoStikerBit.Filename,
		FotoKelengkapan:        fotoKelengkapan.Filename,
		StatusBarang:           statusBarang,
		Brand:                  brand,
		SnLed:                  snLed,
		TextNotes:              textNotes,
	}

	response = controller.prestagingDigitalSignageService.PostPrestaging(requestMap, request)
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

func (controller *prestagingDigitalSignageController) ApprovePrestagingDigitalSignage(ctx *gin.Context) {
	var request prestagingDigitalSignageRequest.ApprovePrestaging
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
	response = controller.prestagingDigitalSignageService.ApprovePrestaging(request)
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

func (controller *prestagingDigitalSignageController) RejectPrestagingDigitalSignage(ctx *gin.Context) {
	var request prestagingDigitalSignageRequest.RejectPrestaging
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
	response = controller.prestagingDigitalSignageService.RejectPrestaging(request)
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

func (controller *prestagingDigitalSignageController) ReuploadPrestagingDigitalSignage(ctx *gin.Context) {
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
	fotoLedFull, _ := ctx.FormFile("fotoLedFull")
	fotoSnLed, _ := ctx.FormFile("fotoSnLed")
	fotoTestDeadPixelPutih, _ := ctx.FormFile("fotoTestDeadPixelPutih")
	fotoTestDeadPixelBiru, _ := ctx.FormFile("fotoTestDeadPixelBiru")
	fotoTestDeadPixelMerah, _ := ctx.FormFile("fotoTestDeadPixelMerah")
	fotoStikerBit, _ := ctx.FormFile("fotoStikerBit")
	fotoTestDeadPixelHitam, _ := ctx.FormFile("fotoTestDeadPixelHitam")
	fotoKelengkapan, _ := ctx.FormFile("fotoKelengkapan")
	fotoTestDeadPixelHijau, _ := ctx.FormFile("fotoTestDeadPixelHijau")
	brand := ctx.PostForm("brand")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoLedFull.Filename] = fotoLedFull
	requestMap[fotoSnLed.Filename] = fotoSnLed
	requestMap[fotoTestDeadPixelPutih.Filename] = fotoTestDeadPixelPutih
	requestMap[fotoTestDeadPixelBiru.Filename] = fotoTestDeadPixelBiru
	requestMap[fotoTestDeadPixelMerah.Filename] = fotoTestDeadPixelMerah
	requestMap[fotoStikerBit.Filename] = fotoStikerBit
	requestMap[fotoTestDeadPixelHitam.Filename] = fotoTestDeadPixelHitam
	requestMap[fotoKelengkapan.Filename] = fotoKelengkapan
	requestMap[fotoTestDeadPixelHijau.Filename] = fotoTestDeadPixelHijau

	request := prestagingDigitalSignageRequest.PostPrestaging{
		IdUploader:             idUploader,
		Uploader:               uploader,
		Sn:                     sn,
		ProjectName:            projectName,
		FotoLedFull:            fotoLedFull.Filename,
		FotoSnLed:              fotoSnLed.Filename,
		FotoTestDeadPixelPutih: fotoTestDeadPixelPutih.Filename,
		FotoTestDeadPixelBiru:  fotoTestDeadPixelBiru.Filename,
		FotoTestDeadPixelMerah: fotoTestDeadPixelMerah.Filename,
		FotoTestDeadPixelHijau: fotoTestDeadPixelHijau.Filename,
		FotoTestDeadPixelHitam: fotoTestDeadPixelHitam.Filename,
		FotoStikerBit:          fotoStikerBit.Filename,
		FotoKelengkapan:        fotoKelengkapan.Filename,
		Brand:                  brand,
		TextNotes:              textNotes,
	}

	response = controller.prestagingDigitalSignageService.ReuploadPrestaging(requestMap, request)
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

func (controller *prestagingDigitalSignageController) AllSubmittedDataPrestagingDigitalSignage(ctx *gin.Context) {
	response := controller.prestagingDigitalSignageService.AllSubmittedData()
	ctx.JSON(response.HttpCode, response)
}

func (controller *prestagingDigitalSignageController) GetSubmittedDataPrestagingDigitalSignageBySn(ctx *gin.Context) {
	var request prestagingDigitalSignageRequest.FindBySn
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
	response = controller.prestagingDigitalSignageService.GetSubmittedDataBySn(request)
	ctx.JSON(response.HttpCode, response)
}

func (controller *prestagingDigitalSignageController) GetRejectedDataPrestagingDigitalSignage(ctx *gin.Context) {
	var request prestagingDigitalSignageRequest.FindRejectedData
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
	response = controller.prestagingDigitalSignageService.GetRejectedData(request)
	ctx.JSON(response.HttpCode, response)
}
