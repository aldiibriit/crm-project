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
	AllSubmittedData(ctx *gin.Context)
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
	textNotes := ctx.PostForm("textNotes")
	fotoMesinCrmFull, _ := ctx.FormFile("fotoMesinCrmFull")
	fotoSnMesinCrm, _ := ctx.FormFile("fotoSnMesinCrm")
	fotoCameraAtas, _ := ctx.FormFile("fotoCameraAtas")
	fotoCameraCashOut, _ := ctx.FormFile("fotoCameraCashOut")
	fotoSystemInformationCu, _ := ctx.FormFile("fotoSystemInformationCu")
	fotoKapasitasHardisk, _ := ctx.FormFile("fotoKapasitasHardisk")
	fotoKunciCrm, _ := ctx.FormFile("fotoKunciCrm")
	fotoClr, _ := ctx.FormFile("fotoClr")
	fotoPortPc, _ := ctx.FormFile("fotoPortPc")
	fotoStikerBit, _ := ctx.FormFile("fotoStikerBit")
	fotoStrukErrorLogTest, _ := ctx.FormFile("fotoStrukErrorLogTest")
	fotoCeklist, _ := ctx.FormFile("fotoCeklist")
	snCpu := ctx.PostForm("snCpu")
	snClr := ctx.PostForm("snClr")
	snReceiptPrinter := ctx.PostForm("snReceiptPrinter")
	snUr := ctx.PostForm("snUr")
	snBv := ctx.PostForm("snBv")
	statusDeadPixelMonitor := ctx.PostForm("statusDeadPixelMonitor")
	snMonitor := ctx.PostForm("snMonitor")
	brand := ctx.PostForm("brand")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoMesinCrmFull.Filename] = fotoMesinCrmFull

	request := prestagingCRMRequest.PostPrestaging{
		IdUploader:              idUploader,
		Uploader:                uploader,
		Sn:                      sn,
		ProjectName:             projectName,
		FotoMesinCrmFull:        fotoMesinCrmFull.Filename,
		FotoSnMesinCrm:          fotoSnMesinCrm.Filename,
		FotoCameraAtas:          fotoCameraAtas.Filename,
		FotoCameraCashOut:       fotoCameraCashOut.Filename,
		FotoSystemInformationCu: fotoSystemInformationCu.Filename,
		FotoKapasitasHardisk:    fotoKapasitasHardisk.Filename,
		FotoKunciCrm:            fotoKunciCrm.Filename,
		FotoClr:                 fotoClr.Filename,
		FotoPortPc:              fotoPortPc.Filename,
		FotoStikerBit:           fotoStikerBit.Filename,
		FotoStrukErrorLogTest:   fotoStrukErrorLogTest.Filename,
		FotoCeklist:             fotoCeklist.Filename,
		SnCpu:                   snCpu,
		SnClr:                   snClr,
		SnReceiptPrinter:        snReceiptPrinter,
		SnUr:                    snUr,
		SnBv:                    snBv,
		SnMonitor:               snMonitor,
		StatusDeadPixelMonitor:  statusDeadPixelMonitor,
		Brand:                   brand,
		TextNotes:               textNotes,
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
	textNotes := ctx.PostForm("textNotes")
	fotoMesinCrmFull, _ := ctx.FormFile("fotoMesinCrmFull")
	fotoSnMesinCrm, _ := ctx.FormFile("fotoSnMesinCrm")
	fotoCameraAtas, _ := ctx.FormFile("fotoCameraAtas")
	fotoCameraCashOut, _ := ctx.FormFile("fotoCameraCashOut")
	fotoSystemInformationCu, _ := ctx.FormFile("fotoSystemInformationCu")
	fotoKapasitasHardisk, _ := ctx.FormFile("fotoKapasitasHardisk")
	fotoKunciCrm, _ := ctx.FormFile("fotoKunciCrm")
	fotoClr, _ := ctx.FormFile("fotoClr")
	fotoPortPc, _ := ctx.FormFile("fotoPortPc")
	fotoStikerBit, _ := ctx.FormFile("fotoStikerBit")
	fotoStrukErrorLogTest, _ := ctx.FormFile("fotoStrukErrorLogTest")
	fotoCeklist, _ := ctx.FormFile("fotoCeklist")
	snCpu := ctx.PostForm("snCpu")
	snClr := ctx.PostForm("snClr")
	snReceiptPrinter := ctx.PostForm("snReceiptPrinter")
	snUr := ctx.PostForm("snUr")
	snBv := ctx.PostForm("snBv")
	statusDeadPixelMonitor := ctx.PostForm("statusDeadPixelMonitor")
	snMonitor := ctx.PostForm("snMonitor")
	brand := ctx.PostForm("brand")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoMesinCrmFull.Filename] = fotoMesinCrmFull

	request := prestagingCRMRequest.PostPrestaging{
		IdUploader:              idUploader,
		Uploader:                uploader,
		Sn:                      sn,
		ProjectName:             projectName,
		FotoMesinCrmFull:        fotoMesinCrmFull.Filename,
		FotoSnMesinCrm:          fotoSnMesinCrm.Filename,
		FotoCameraAtas:          fotoCameraAtas.Filename,
		FotoCameraCashOut:       fotoCameraCashOut.Filename,
		FotoSystemInformationCu: fotoSystemInformationCu.Filename,
		FotoKapasitasHardisk:    fotoKapasitasHardisk.Filename,
		FotoKunciCrm:            fotoKunciCrm.Filename,
		FotoClr:                 fotoClr.Filename,
		FotoPortPc:              fotoPortPc.Filename,
		FotoStikerBit:           fotoStikerBit.Filename,
		FotoStrukErrorLogTest:   fotoStrukErrorLogTest.Filename,
		FotoCeklist:             fotoCeklist.Filename,
		SnCpu:                   snCpu,
		SnClr:                   snClr,
		SnReceiptPrinter:        snReceiptPrinter,
		SnUr:                    snUr,
		SnBv:                    snBv,
		SnMonitor:               snMonitor,
		StatusDeadPixelMonitor:  statusDeadPixelMonitor,
		Brand:                   brand,
		TextNotes:               textNotes,
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

func (controller *prestagingCRMController) AllSubmittedData(ctx *gin.Context) {
	response := controller.prestagingCRMService.AllSubmittedData()
	ctx.JSON(response.HttpCode, response)
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
