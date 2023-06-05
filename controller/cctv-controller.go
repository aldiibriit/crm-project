package controller

import (
	"errors"
	CCTVRequestDTO "go-api/dto/request/cctv"
	"go-api/dto/response"
	"go-api/helper"
	"go-api/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CCTVController interface {
	GetAll(ctx *gin.Context)
	FindBySN(ctx *gin.Context)
}

type cctvController struct {
	cctvService service.CCTVService
	jwtService  service.JWTService
}

func NewCCTVController(cctvServ service.CCTVService, jwtServ service.JWTService) CCTVController {
	return &cctvController{
		cctvService: cctvServ,
		jwtService:  jwtServ,
	}
}

func (controller *cctvController) GetAll(ctx *gin.Context) {
	response := controller.cctvService.GetAll()
	ctx.JSON(response.HttpCode, response)
}

func (controller *cctvController) FindBySN(ctx *gin.Context) {
	var response response.UniversalResponse
	var request CCTVRequestDTO.FindBySNRequest
	var ve validator.ValidationErrors

	err := ctx.ShouldBind(&request)
	if err != nil {
		if errors.As(err, &ve) {
			log.Println(err)
			log.Println(&ve)
			errMessage := helper.CustomValidator(ve)
			response.Data = errMessage
			response.HttpCode = 400
			response.ResponseCode = "99"
			response.ResponseMessage = "Missing parameter"
			ctx.AbortWithStatusJSON(response.HttpCode, response)
			return
		}
	}

	response = controller.cctvService.FindBySN(request)
	ctx.JSON(response.HttpCode, response)
}
