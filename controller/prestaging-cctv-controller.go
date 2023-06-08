package controller

import (
	"errors"
	"go-api/helper"
	"go-api/service"
	"mime/multipart"
	"net/http"
	"os"

	prestagingCctvRequest "go-api/dto/request/prestaging-cctv"
	"go-api/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PrestagingCCTVController interface {
	PostPrestagingCCTV(ctx *gin.Context)
	ApprovePrestagingCCTV(ctx *gin.Context)
	RejectPrestagingCCTV(ctx *gin.Context)
	ReuploadPrestagingCCTV(ctx *gin.Context)
	AllSubmittedDataPrestagingCCTV(ctx *gin.Context)
	GetSubmittedDataPrestagingCCTVBySn(ctx *gin.Context)
	GetRejectedDataPrestagingCCTV(ctx *gin.Context)
}

type prestagingCCTVController struct {
	prestagingCCTVService service.PrestagingCCTVService
	jwtService            service.JWTService
}

func NewPrestagingCCTVController(prestagingCCTVServ service.PrestagingCCTVService, jwtServ service.JWTService) PrestagingCCTVController {
	return &prestagingCCTVController{
		prestagingCCTVService: prestagingCCTVServ,
		jwtService:            jwtServ,
	}
}

func (controller *prestagingCCTVController) PostPrestagingCCTV(ctx *gin.Context) {
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
	fotoSnNvr, _ := ctx.FormFile("fotoSnNvr")
	fotoSnCamera, _ := ctx.FormFile("fotoSnCamera")
	fotoKelengkapanNvr, _ := ctx.FormFile("fotoKelengkapanNvr")
	fotoKelengkapanCamera, _ := ctx.FormFile("fotoKelengkapanCamera")
	fotoInstallasiInputDanOutputPortNvr, _ := ctx.FormFile("fotoInstallasiInputDanOutputPortNvr")
	fotoTampilanCamera, _ := ctx.FormFile("fotoTampilanCamera")
	fotoSettingResolusi, _ := ctx.FormFile("fotoSettingResolusi")
	fotoSettingMotion, _ := ctx.FormFile("fotoSettingMotion")
	fotoStorage, _ := ctx.FormFile("fotoStorage")
	fotoCheckBackup, _ := ctx.FormFile("fotoCheckBackup")
	fotoStikerBit, _ := ctx.FormFile("fotoStikerBit")
	snNvr := ctx.PostForm("snNvr")
	snCamera := ctx.PostForm("snCamera")
	userNewNvr := ctx.PostForm("userNewNvr")
	passwordNewNvr := ctx.PostForm("passwordNewNvr")
	statusBarang := ctx.PostForm("statusBarang")
	brand := ctx.PostForm("brand")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoNvrDanCamera.Filename] = fotoNvrDanCamera
	requestMap[fotoSnNvr.Filename] = fotoSnNvr
	requestMap[fotoSnCamera.Filename] = fotoSnCamera
	requestMap[fotoKelengkapanNvr.Filename] = fotoKelengkapanNvr
	requestMap[fotoKelengkapanCamera.Filename] = fotoKelengkapanCamera
	requestMap[fotoInstallasiInputDanOutputPortNvr.Filename] = fotoInstallasiInputDanOutputPortNvr
	requestMap[fotoTampilanCamera.Filename] = fotoTampilanCamera
	requestMap[fotoTampilanCamera.Filename] = fotoTampilanCamera
	requestMap[fotoSettingResolusi.Filename] = fotoSettingResolusi
	requestMap[fotoSettingMotion.Filename] = fotoSettingMotion
	requestMap[fotoStorage.Filename] = fotoStorage
	requestMap[fotoCheckBackup.Filename] = fotoCheckBackup
	requestMap[fotoStikerBit.Filename] = fotoStikerBit

	request := prestagingCctvRequest.PostPrestaging{
		IdUploader:                          idUploader,
		Uploader:                            uploader,
		Sn:                                  sn,
		ProjectName:                         projectName,
		FotoNvrDanCamera:                    fotoNvrDanCamera.Filename,
		FotoSnNvr:                           fotoSnNvr.Filename,
		FotoSnCamera:                        fotoSnCamera.Filename,
		FotoKelengkapanNvr:                  fotoKelengkapanNvr.Filename,
		FotoKelengkapanCamera:               fotoKelengkapanCamera.Filename,
		FotoInstallasiInputDanOutputPortNvr: fotoInstallasiInputDanOutputPortNvr.Filename,
		FotoTampilanCamera:                  fotoTampilanCamera.Filename,
		FotoSettingResolusi:                 fotoSettingResolusi.Filename,
		FotoSettingMotion:                   fotoSettingMotion.Filename,
		FotoStorage:                         fotoStorage.Filename,
		FotoCheckBackup:                     fotoCheckBackup.Filename,
		FotoStikerBit:                       fotoStikerBit.Filename,
		StatusBarang:                        statusBarang,
		UserNewNvr:                          userNewNvr,
		PasswordNewNvr:                      passwordNewNvr,
		Brand:                               brand,
		SnNvr:                               snNvr,
		SnCamera:                            snCamera,
		TextNotes:                           textNotes,
	}

	response = controller.prestagingCCTVService.PostPrestaging(requestMap, request)
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

func (controller *prestagingCCTVController) ApprovePrestagingCCTV(ctx *gin.Context) {
	var request prestagingCctvRequest.ApprovePrestaging
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
	response = controller.prestagingCCTVService.ApprovePrestaging(request)
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

func (controller *prestagingCCTVController) RejectPrestagingCCTV(ctx *gin.Context) {
	var request prestagingCctvRequest.RejectPrestaging
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
	response = controller.prestagingCCTVService.RejectPrestaging(request)
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

func (controller *prestagingCCTVController) ReuploadPrestagingCCTV(ctx *gin.Context) {
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
	fotoSnNvr, _ := ctx.FormFile("fotoSnNvr")
	fotoSnCamera, _ := ctx.FormFile("fotoSnCamera")
	fotoKelengkapanNvr, _ := ctx.FormFile("fotoKelengkapanNvr")
	fotoKelengkapanCamera, _ := ctx.FormFile("fotoKelengkapanCamera")
	fotoInstallasiInputDanOutputPortNvr, _ := ctx.FormFile("fotoInstallasiInputDanOutputPortNvr")
	fotoTampilanCamera, _ := ctx.FormFile("fotoTampilanCamera")
	fotoSettingResolusi, _ := ctx.FormFile("fotoSettingResolusi")
	fotoSettingMotion, _ := ctx.FormFile("fotoSettingMotion")
	fotoStorage, _ := ctx.FormFile("fotoStorage")
	fotoCheckBackup, _ := ctx.FormFile("fotoCheckBackup")
	fotoStikerBit, _ := ctx.FormFile("fotoStikerBit")
	brand := ctx.PostForm("brand")
	// log.Println(fotoSnCrm.Filename)
	requestMap := make(map[string]*multipart.FileHeader)
	requestMap[fotoNvrDanCamera.Filename] = fotoNvrDanCamera
	requestMap[fotoSnNvr.Filename] = fotoSnNvr
	requestMap[fotoSnCamera.Filename] = fotoSnCamera
	requestMap[fotoKelengkapanNvr.Filename] = fotoKelengkapanNvr
	requestMap[fotoKelengkapanCamera.Filename] = fotoKelengkapanCamera
	requestMap[fotoInstallasiInputDanOutputPortNvr.Filename] = fotoInstallasiInputDanOutputPortNvr
	requestMap[fotoTampilanCamera.Filename] = fotoTampilanCamera
	requestMap[fotoTampilanCamera.Filename] = fotoTampilanCamera
	requestMap[fotoSettingResolusi.Filename] = fotoSettingResolusi
	requestMap[fotoSettingMotion.Filename] = fotoSettingMotion
	requestMap[fotoStorage.Filename] = fotoStorage
	requestMap[fotoCheckBackup.Filename] = fotoCheckBackup
	requestMap[fotoStikerBit.Filename] = fotoStikerBit

	request := prestagingCctvRequest.PostPrestaging{
		IdUploader:                          idUploader,
		Uploader:                            uploader,
		Sn:                                  sn,
		ProjectName:                         projectName,
		FotoNvrDanCamera:                    fotoNvrDanCamera.Filename,
		FotoSnNvr:                           fotoSnNvr.Filename,
		FotoSnCamera:                        fotoSnCamera.Filename,
		FotoKelengkapanNvr:                  fotoKelengkapanNvr.Filename,
		FotoKelengkapanCamera:               fotoKelengkapanCamera.Filename,
		FotoInstallasiInputDanOutputPortNvr: fotoInstallasiInputDanOutputPortNvr.Filename,
		FotoTampilanCamera:                  fotoTampilanCamera.Filename,
		FotoSettingResolusi:                 fotoSettingResolusi.Filename,
		FotoSettingMotion:                   fotoSettingMotion.Filename,
		FotoStorage:                         fotoStorage.Filename,
		FotoCheckBackup:                     fotoCheckBackup.Filename,
		FotoStikerBit:                       fotoStikerBit.Filename,
		Brand:                               brand,
		TextNotes:                           textNotes,
	}

	response = controller.prestagingCCTVService.ReuploadPrestaging(requestMap, request)
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

func (controller *prestagingCCTVController) AllSubmittedDataPrestagingCCTV(ctx *gin.Context) {
	response := controller.prestagingCCTVService.AllSubmittedData()
	ctx.JSON(response.HttpCode, response)
}

func (controller *prestagingCCTVController) GetSubmittedDataPrestagingCCTVBySn(ctx *gin.Context) {
	var request prestagingCctvRequest.FindBySn
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
	response = controller.prestagingCCTVService.GetSubmittedDataBySn(request)
	ctx.JSON(response.HttpCode, response)
}

func (controller *prestagingCCTVController) GetRejectedDataPrestagingCCTV(ctx *gin.Context) {
	var request prestagingCctvRequest.FindRejectedData
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
	response = controller.prestagingCCTVService.GetRejectedData(request)
	ctx.JSON(response.HttpCode, response)
}
