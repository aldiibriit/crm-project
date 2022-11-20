package controller

import (
	"encoding/base64"
	"fmt"
	"go-api/helper"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

type TestController interface {
	DecryptRequest(context *gin.Context)
	TestEmail(context *gin.Context)
}

type testController struct {
	emailservice service.EmailService
}

func NewTestController(emailServ service.EmailService) TestController {
	return &testController{
		emailservice: emailServ,
	}
}

func (x *testController) DecryptRequest(context *gin.Context) {
	r := context.PostForm("request")
	x2, _ := helper.RsaEncryptFEToBE([]byte(r))
	fmt.Println(string(x2))
	decodedX, _ := base64.StdEncoding.DecodeString(x2)
	plain, _ := helper.RsaDecryptFromFEInBE(decodedX)
	fmt.Println(plain)
}

func (x *testController) TestEmail(context *gin.Context) {
	// x.emailservice.SendMessage("")
	context.JSON(200, gin.H{
		"testing": true,
	})
}
