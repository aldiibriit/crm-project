package controller

import (
	"errors"
	"go-api/helper"
	"go-api/service"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	prestagingCRMRequest "go-api/dto/request/prestaging-crm"
	"go-api/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PrestagingCRMController interface {
	PostPrestaging(ctx *gin.Context)
	ApprovePrestaging(ctx *gin.Context)
	RejectPrestaging(ctx *gin.Context)
	ReuploadPrestaging(ctx *gin.Context)
	GetTotalDataEachStatus(ctx *gin.Context)
	PostPrestagingV2(ctx *gin.Context)
}

type prestagingCRMController struct {
	prestagingCRMService service.PrestagingCRMService
	jwtService           service.JWTService
}

func NewPrestagingCRMController(prestagingCRMServ service.PrestagingCRMService, jwtServ service.JWTService) PrestagingCRMController {
	return &prestagingCRMController{
		prestagingCRMService: prestagingCRMServ,
		jwtService:           jwtServ,
	}
}

func (controller *prestagingCRMController) PostPrestaging(ctx *gin.Context) {
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
	tanggalPrestaging := ctx.PostForm("tanggalPrestaging")
	textNotes := ctx.PostForm("textNotes")
	fotoSnCrm, _ := ctx.FormFile("fotoSnCrm")
	fotoLembarKelengkapanCrm, _ := ctx.FormFile("fotoLembarKelengkapanCrm")
	fotoCheckKameraAtas, _ := ctx.FormFile("fotoCheckKameraAtas")
	fotoCheckKameraSamping, _ := ctx.FormFile("fotoCheckKameraSamping")
	fotoKunciCrm, _ := ctx.FormFile("fotoKunciCrm")
	fotoStrukErrorLog, _ := ctx.FormFile("fotoStrukErrorLog")
	fotoContactlessReader, _ := ctx.FormFile("fotoContactlessReader")
	fotoKomponenPc, _ := ctx.FormFile("fotoKomponenPc")
	fotoStikerBriit, _ := ctx.FormFile("fotoStikerBriit")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoSnCrm.Filename] = fotoSnCrm
	requestMap[fotoLembarKelengkapanCrm.Filename] = fotoLembarKelengkapanCrm
	requestMap[fotoCheckKameraAtas.Filename] = fotoCheckKameraAtas
	requestMap[fotoCheckKameraSamping.Filename] = fotoCheckKameraSamping
	requestMap[fotoKunciCrm.Filename] = fotoKunciCrm
	requestMap[fotoStrukErrorLog.Filename] = fotoStrukErrorLog
	requestMap[fotoContactlessReader.Filename] = fotoContactlessReader
	requestMap[fotoKomponenPc.Filename] = fotoKomponenPc
	requestMap[fotoStikerBriit.Filename] = fotoStikerBriit

	request := prestagingCRMRequest.PostPrestaging{
		IdUploader:               idUploader,
		Uploader:                 uploader,
		Sn:                       sn,
		ProjectName:              projectName,
		FotoSnCrm:                fotoSnCrm.Filename,
		FotoLembarKelengkapanCrm: fotoLembarKelengkapanCrm.Filename,
		FotoCheckKameraAtas:      fotoCheckKameraAtas.Filename,
		FotoCheckKameraSamping:   fotoCheckKameraSamping.Filename,
		FotoKunciCrm:             fotoKunciCrm.Filename,
		FotoStrukErrorLog:        fotoStrukErrorLog.Filename,
		FotoContactlessReader:    fotoContactlessReader.Filename,
		FotoKomponenPc:           fotoKomponenPc.Filename,
		FotoStikerBriit:          fotoStikerBriit.Filename,
		TanggalPrestaging:        tanggalPrestaging,
		TextNotes:                textNotes,
	}

	response = controller.prestagingCRMService.PostPrestaging(requestMap, request)
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

func (controller *prestagingCRMController) ApprovePrestaging(ctx *gin.Context) {
	var request prestagingCRMRequest.ApprovePrestaging
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
	response = controller.prestagingCRMService.Approve(request)
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

func (controller *prestagingCRMController) RejectPrestaging(ctx *gin.Context) {
	var request prestagingCRMRequest.RejectPrestaging
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
	response = controller.prestagingCRMService.Reject(request)
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

func (controller *prestagingCRMController) ReuploadPrestaging(ctx *gin.Context) {
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
	tanggalPrestaging := ctx.PostForm("tanggalPrestaging")
	textNotes := ctx.PostForm("textNotes")
	fotoSnCrm, _ := ctx.FormFile("fotoSnCrm")
	fotoLembarKelengkapanCrm, _ := ctx.FormFile("fotoLembarKelengkapanCrm")
	fotoCheckKameraAtas, _ := ctx.FormFile("fotoCheckKameraAtas")
	fotoCheckKameraSamping, _ := ctx.FormFile("fotoCheckKameraSamping")
	fotoKunciCrm, _ := ctx.FormFile("fotoKunciCrm")
	fotoStrukErrorLog, _ := ctx.FormFile("fotoStrukErrorLog")
	fotoContactlessReader, _ := ctx.FormFile("fotoContactlessReader")
	fotoKomponenPc, _ := ctx.FormFile("fotoKomponenPc")
	fotoStikerBriit, _ := ctx.FormFile("fotoStikerBriit")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoSnCrm.Filename] = fotoSnCrm
	requestMap[fotoLembarKelengkapanCrm.Filename] = fotoLembarKelengkapanCrm
	requestMap[fotoCheckKameraAtas.Filename] = fotoCheckKameraAtas
	requestMap[fotoCheckKameraSamping.Filename] = fotoCheckKameraSamping
	requestMap[fotoKunciCrm.Filename] = fotoKunciCrm
	requestMap[fotoStrukErrorLog.Filename] = fotoStrukErrorLog
	requestMap[fotoContactlessReader.Filename] = fotoContactlessReader
	requestMap[fotoKomponenPc.Filename] = fotoKomponenPc
	requestMap[fotoStikerBriit.Filename] = fotoStikerBriit

	request := prestagingCRMRequest.PostPrestaging{
		IdUploader:               idUploader,
		Uploader:                 uploader,
		Sn:                       sn,
		ProjectName:              projectName,
		FotoSnCrm:                fotoSnCrm.Filename,
		FotoLembarKelengkapanCrm: fotoLembarKelengkapanCrm.Filename,
		FotoCheckKameraAtas:      fotoCheckKameraAtas.Filename,
		FotoCheckKameraSamping:   fotoCheckKameraSamping.Filename,
		FotoKunciCrm:             fotoKunciCrm.Filename,
		FotoStrukErrorLog:        fotoStrukErrorLog.Filename,
		FotoContactlessReader:    fotoContactlessReader.Filename,
		FotoKomponenPc:           fotoKomponenPc.Filename,
		FotoStikerBriit:          fotoStikerBriit.Filename,
		TanggalPrestaging:        tanggalPrestaging,
		TextNotes:                textNotes,
	}

	response = controller.prestagingCRMService.ReuploadPrestaging(requestMap, request)
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

func (controller *prestagingCRMController) GetTotalDataEachStatus(ctx *gin.Context) {

}

func (controller *prestagingCRMController) PostPrestagingV2(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Println(err.Error())
	}

	files := form.File["files"]

	for _, file := range files {
		log.Println(file.Filename)
	}

	response := controller.prestagingCRMService.PostPrestagingV2(files)

	ctx.JSON(response.HttpCode, response)
}
