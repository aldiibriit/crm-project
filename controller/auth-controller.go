package controller

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"go-api/dto"
	"go-api/dto/request/authRequestDTO"
	"go-api/entity"
	"go-api/helper"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// AuthController interface is a contract what this controller can do
type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	RegisterSales(ctx *gin.Context)
	CreateToken(ctx *gin.Context)
	ActivateUser(ctx *gin.Context)
	PasswordConfirmation(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {

	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// decrypt request from FE
	plainRequest := deserializeLoginRequest(loginDTO)

	authResult := c.authService.VerifyCredential(plainRequest.Email, plainRequest.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken
		response := helper.BuildResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.BuildResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}

func (c *authController) RegisterSales(ctx *gin.Context) {
	var registerSalesDTO dto.RegisterSalesDTO
	errDTO := ctx.ShouldBind(&registerSalesDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// if !c.authService.IsDuplicateEmail(registerSalesDTO.EmailSales) {
	// 	response := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
	// 	ctx.JSON(http.StatusConflict, response)
	// } else {
	createdUser := c.authService.CreateUserSales(registerSalesDTO)
	ctx.JSON(createdUser.HttpCode, createdUser)
}

func (c *authController) ActivateUser(ctx *gin.Context) {
	var activateRequestDTO authRequestDTO.ActivateRequestDTO
	errDTO := ctx.ShouldBind(&activateRequestDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := c.authService.ActivateUser(activateRequestDTO)
	ctx.JSON(res.HttpCode, res)
}

func (c *authController) CreateToken(ctx *gin.Context) {
	log.Println("AuthController|CreateToken")
	var authRequestDTO authRequestDTO.AuthRequest

	errDTO := ctx.ShouldBind(&authRequestDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authRequestDecrypted := deserializeCreateTokenRequest(authRequestDTO)
	authRequestDTO.Email = authRequestDecrypted.Email
	x := rand.Intn(10)
	var xChangeId string = strconv.Itoa(x)
	// var metadata map[string]interface{}
	token := c.jwtService.GenerateToken2(xChangeId, authRequestDTO)

	ctx.JSON(200, gin.H{
		"decryptedEmail":  authRequestDTO.Email,
		"applicationName": authRequestDTO.ApplicationName,
		"token":           token,
	})

}

func (c *authController) PasswordConfirmation(ctx *gin.Context) {
	var passwordConfirmationDTO dto.PasswordConfirmationDTO

	errDTO := ctx.ShouldBind(&passwordConfirmationDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := c.authService.PasswordConfirmation(passwordConfirmationDTO)
	ctx.JSON(res.HttpCode, res)
}

// returning plain text from encrypted request
func deserializeLoginRequest(request interface{}) dto.LoginDTO {
	loginDTO := request.(dto.LoginDTO)

	cipheTextEmail, _ := base64.StdEncoding.DecodeString(loginDTO.Email)
	cipheTextPassword, _ := base64.StdEncoding.DecodeString(loginDTO.Password)
	plainTextEmail, _ := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmail))
	plainTextPassword, _ := helper.RsaDecryptFromFEInBE([]byte(cipheTextPassword))

	var result dto.LoginDTO

	result.Email = plainTextEmail
	result.Password = plainTextPassword

	return result
}

func deserializeCreateTokenRequest(request interface{}) authRequestDTO.AuthRequest {
	authRequest := request.(authRequestDTO.AuthRequest)

	cipheTextEmail, _ := base64.StdEncoding.DecodeString(authRequest.Email)
	plainTextEmail, _ := helper.RsaDecryptFromFEInBE([]byte(cipheTextEmail))

	var result authRequestDTO.AuthRequest

	result.Email = plainTextEmail

	return result
}
