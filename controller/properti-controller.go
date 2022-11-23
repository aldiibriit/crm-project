package controller

import (
	propertiDTO "go-api/dto/properti"
	"go-api/helper"
	"go-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PropertiController interface {
	AdvancedFilter(context *gin.Context)
	LandingPage(context *gin.Context)
}

type propertiController struct {
	propertiService service.PropertiService
}

// NewPropertiController create a new instances of BoookController
func NewPropertiController(propertiServ service.PropertiService) PropertiController {
	return &propertiController{
		propertiService: propertiServ,
	}
}

func (c *propertiController) AdvancedFilter(context *gin.Context) {
	var advancedFilterDTO propertiDTO.AdvancedFilterDTO
	bearerToken := context.GetHeader("Authorization")
	errDTO := context.ShouldBind(&advancedFilterDTO)

	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result := c.propertiService.AdvancedFilter(advancedFilterDTO, bearerToken)
	context.JSON(200, result)
}

func (c *propertiController) LandingPage(ctx *gin.Context) {
	result := c.propertiService.LandingPage()
	ctx.JSON(result.HttpCode, result)
}
