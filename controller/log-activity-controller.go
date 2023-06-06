package controller

import (
	"errors"
	LogActivityRequest "go-api/dto/request/log-activity"
	"go-api/dto/response"
	"go-api/helper"
	"go-api/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LogActivityController interface {
	GetTimeLine(ctx *gin.Context)
}

type logActivityController struct {
	logActivityService service.LogActivityService
}

func NewLogActivityController(logActivityServ service.LogActivityService) LogActivityController {
	return &logActivityController{
		logActivityService: logActivityServ,
	}
}

func (controller *logActivityController) GetTimeLine(ctx *gin.Context) {
	var request LogActivityRequest.GetTimeLine

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
	successResponse := controller.logActivityService.GetTimeLine(request)
	ctx.JSON(successResponse.HttpCode, successResponse)
}
