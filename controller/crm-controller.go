package controller

import (
	"go-api/service"

	"github.com/gin-gonic/gin"
)

type CRMController interface {
	GetAll(ctx *gin.Context)
}

type crmController struct {
	crmService service.CRMService
	jwtService service.JWTService
}

func NewCRMController(crmServ service.CRMService, jwtServ service.JWTService) CRMController {
	return &crmController{
		crmService: crmServ,
		jwtService: jwtServ,
	}
}

func (controller *crmController) GetAll(ctx *gin.Context) {

}
