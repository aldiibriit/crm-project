package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-api/helper"
	"go-api/service"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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
	url := "https://api.privateopen.sandbox.bri.co.id/gateway/Oauth2/1.0/accessToken"
	method := "POST"

	payload := strings.NewReader("client_id=b78f2b47-b258-49cb-9a0d-cab6f4d5f80a&client_secret=eff57260-b4c2-4b5d-b04a-b59b2a9f51d1&grant_type=client_credentials&scope=apiKPRBrispot")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "TS015ba640=01d4c679fed33239583546a1c96b3f8fa7de4d7e56256a5d3c734980c774b4ebceb059cdd06f47912fa72d0836bf7361dcde55bea8")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var responseCreateToken map[string]interface{}
	accessToken := fmt.Sprintf("%v", responseCreateToken["access_token"])
	json.Unmarshal([]byte(string(body)), &responseCreateToken)
	timestamp := time.Now().Format(time.RFC3339)
	signature := helper.HMACBuilder("/gateway/apiKPRBrispot/1.0/insertPrakarsaMortgage", "POST", "Bearer "+accessToken, "", "eff57260-b4c2-4b5d-b04a-b59b2a9f51d1", timestamp)
	context.JSON(200, gin.H{
		"testing":   responseCreateToken["access_token"],
		"signature": signature,
	})
}
