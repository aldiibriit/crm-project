package controller

import (
	"errors"
	"go-api/helper"
	"go-api/service"
	"mime/multipart"
	"net/http"
	"os"

	prestagingUPSService "go-api/dto/request/prestaging-ups"
	"go-api/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PrestagingUPSController interface {
	PostPrestagingUPS(ctx *gin.Context)
	ApprovePrestagingUPS(ctx *gin.Context)
	RejectPrestagingUPS(ctx *gin.Context)
	ReuploadPrestagingUPS(ctx *gin.Context)
	AllSubmittedDataPrestagingUPS(ctx *gin.Context)
	GetSubmittedDataPrestagingUPSBySn(ctx *gin.Context)
	GetRejectedDataPrestagingUPS(ctx *gin.Context)
}

type prestagingUPSController struct {
	prestagingUpsService service.PrestagingUPSService
	jwtService           service.JWTService
}

func NewPrestagingUPSController(prestagingUpsServ service.PrestagingUPSService, jwtServ service.JWTService) PrestagingUPSController {
	return &prestagingUPSController{
		prestagingUpsService: prestagingUpsServ,
		jwtService:           jwtServ,
	}
}

func (controller *prestagingUPSController) PostPrestagingUPS(ctx *gin.Context) {
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
	fotoSnUps, _ := ctx.FormFile("fotoSnUps")
	fotoKapasitasBaterai, _ := ctx.FormFile("fotoKapasitasBaterai")
	fotoTampilanKelistrikanListrikOn, _ := ctx.FormFile("fotoTampilanKelistrikanListrikOn")
	fotoTampilanKelistrikanListrikOff, _ := ctx.FormFile("fotoTampilanKelistrikanListrikOff")
	fotoStikerBit, _ := ctx.FormFile("fotoStikerBit")
	fotoTampakBelakang, _ := ctx.FormFile("fotoTampakBelakang")
	fotoKelengkapan, _ := ctx.FormFile("fotoKelengkapan")
	fotoCeklist, _ := ctx.FormFile("fotoCeklist")
	statusBarang := ctx.PostForm("statusBarang")
	brand := ctx.PostForm("brand")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoUpsFull.Filename] = fotoUpsFull
	requestMap[fotoSnUps.Filename] = fotoSnUps
	requestMap[fotoKapasitasBaterai.Filename] = fotoKapasitasBaterai
	requestMap[fotoTampilanKelistrikanListrikOn.Filename] = fotoTampilanKelistrikanListrikOn
	requestMap[fotoTampilanKelistrikanListrikOff.Filename] = fotoTampilanKelistrikanListrikOff
	requestMap[fotoStikerBit.Filename] = fotoStikerBit
	requestMap[fotoTampakBelakang.Filename] = fotoTampakBelakang
	requestMap[fotoKelengkapan.Filename] = fotoKelengkapan
	requestMap[fotoCeklist.Filename] = fotoCeklist

	request := prestagingUPSService.PostPrestaging{
		IdUploader:                        idUploader,
		Uploader:                          uploader,
		Sn:                                sn,
		ProjectName:                       projectName,
		FotoUpsFull:                       fotoUpsFull.Filename,
		FotoSnUps:                         fotoSnUps.Filename,
		FotoKapasitasBaterai:              fotoKapasitasBaterai.Filename,
		FotoTampilanKelistrikanListrikOn:  fotoTampilanKelistrikanListrikOn.Filename,
		FotoTampilanKelistrikanListrikOff: fotoTampilanKelistrikanListrikOff.Filename,
		FotoKelengkapan:                   fotoKelengkapan.Filename,
		FotoCeklist:                       fotoCeklist.Filename,
		FotoTampakBelakang:                fotoTampakBelakang.Filename,
		FotoStikerBit:                     fotoStikerBit.Filename,
		StatusBarang:                      statusBarang,
		Brand:                             brand,
		TextNotes:                         textNotes,
	}

	response = controller.prestagingUpsService.PostPrestaging(requestMap, request)
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

func (controller *prestagingUPSController) ApprovePrestagingUPS(ctx *gin.Context) {
	var request prestagingUPSService.ApprovePrestaging
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
	response = controller.prestagingUpsService.ApprovePrestaging(request)
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

func (controller *prestagingUPSController) RejectPrestagingUPS(ctx *gin.Context) {
	var request prestagingUPSService.RejectPrestaging
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
	response = controller.prestagingUpsService.RejectPrestaging(request)
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

func (controller *prestagingUPSController) ReuploadPrestagingUPS(ctx *gin.Context) {
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
	fotoSnUps, _ := ctx.FormFile("fotoSnUps")
	fotoKapasitasBaterai, _ := ctx.FormFile("fotoKapasitasBaterai")
	fotoTampilanKelistrikanListrikOn, _ := ctx.FormFile("fotoTampilanKelistrikanListrikOn")
	fotoTampilanKelistrikanListrikOff, _ := ctx.FormFile("fotoTampilanKelistrikanListrikOff")
	fotoStikerBit, _ := ctx.FormFile("fotoStikerBit")
	fotoTampakBelakang, _ := ctx.FormFile("fotoTampakBelakang")
	fotoKelengkapan, _ := ctx.FormFile("fotoKelengkapan")
	fotoCeklist, _ := ctx.FormFile("fotoCeklist")
	brand := ctx.PostForm("brand")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoUpsFull.Filename] = fotoUpsFull
	requestMap[fotoSnUps.Filename] = fotoSnUps
	requestMap[fotoKapasitasBaterai.Filename] = fotoKapasitasBaterai
	requestMap[fotoTampilanKelistrikanListrikOn.Filename] = fotoTampilanKelistrikanListrikOn
	requestMap[fotoTampilanKelistrikanListrikOff.Filename] = fotoTampilanKelistrikanListrikOff
	requestMap[fotoStikerBit.Filename] = fotoStikerBit
	requestMap[fotoTampakBelakang.Filename] = fotoTampakBelakang
	requestMap[fotoKelengkapan.Filename] = fotoKelengkapan
	requestMap[fotoCeklist.Filename] = fotoCeklist

	request := prestagingUPSService.PostPrestaging{
		IdUploader:                        idUploader,
		Uploader:                          uploader,
		Sn:                                sn,
		ProjectName:                       projectName,
		FotoUpsFull:                       fotoUpsFull.Filename,
		FotoSnUps:                         fotoSnUps.Filename,
		FotoKapasitasBaterai:              fotoKapasitasBaterai.Filename,
		FotoTampilanKelistrikanListrikOn:  fotoTampilanKelistrikanListrikOn.Filename,
		FotoTampilanKelistrikanListrikOff: fotoTampilanKelistrikanListrikOff.Filename,
		FotoKelengkapan:                   fotoKelengkapan.Filename,
		FotoCeklist:                       fotoCeklist.Filename,
		FotoTampakBelakang:                fotoTampakBelakang.Filename,
		FotoStikerBit:                     fotoStikerBit.Filename,
		Brand:                             brand,
		TextNotes:                         textNotes,
	}

	response = controller.prestagingUpsService.ReuploadPrestaging(requestMap, request)
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

func (controller *prestagingUPSController) AllSubmittedDataPrestagingUPS(ctx *gin.Context) {
	response := controller.prestagingUpsService.AllSubmittedData()
	ctx.JSON(response.HttpCode, response)
}

func (controller *prestagingUPSController) GetSubmittedDataPrestagingUPSBySn(ctx *gin.Context) {
	var request prestagingUPSService.FindBySn
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
	response = controller.prestagingUpsService.GetSubmittedDataBySn(request)
	ctx.JSON(response.HttpCode, response)
}

func (controller *prestagingUPSController) GetRejectedDataPrestagingUPS(ctx *gin.Context) {
	var request prestagingUPSService.FindRejectedData
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
	response = controller.prestagingUpsService.GetRejectedData(request)
	ctx.JSON(response.HttpCode, response)
}
