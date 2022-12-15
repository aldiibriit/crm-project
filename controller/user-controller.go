package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"go-api/dto"
	"go-api/dto/request/userRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/helper"
	"go-api/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// UserController is a ....
type UserController interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
	GetDeveloper(ctx *gin.Context)
	GetUserReferral(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

// NewUserController is creating anew instance of UserControlller
func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) Update(context *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errDTO := context.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	userUpdateDTO.ID = id
	u := c.userService.Update(userUpdateDTO)
	res := helper.BuildResponse(true, "OK!", u)
	context.JSON(http.StatusOK, res)
}

func (c *userController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user := c.userService.Profile(id)
	res := helper.BuildResponse(true, "OK", user)
	context.JSON(http.StatusOK, res)

}

func (c *userController) GetDeveloper(ctx *gin.Context) {
	var request userRequestDTO.ListUserDeveloperRequestDTO
	var response responseDTO.Response
	errDTO := ctx.ShouldBind(&request)
	if errDTO != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	response = c.userService.GetDeveloper(request)

	ctx.JSON(response.HttpCode, response)
}

func (c *userController) GetUserReferral(ctx *gin.Context) {
	var request userRequestDTO.ListUserReferralRequestDTO
	var response responseDTO.Response

	errDTO := ctx.ShouldBind(&request)
	if errDTO != nil {
		response.HttpCode = 400
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = errDTO.Error()
		response.Summary = nil
		ctx.AbortWithStatusJSON(response.HttpCode, response)
		return
	}

	if request.SalesEmail != "" {
		decryptedData, err := deserializeGetUserReferralRequest(request)
		if err != nil {
			response.HttpCode = 400
			response.MetadataResponse = nil
			response.ResponseCode = "99"
			response.ResponseData = nil
			response.ResponseDesc = err.Error()
			response.Summary = nil
			ctx.AbortWithStatusJSON(response.HttpCode, response)
			return
		}

		response = c.userService.ListUserReferral(decryptedData)
		ctx.JSON(response.HttpCode, response)
		return
	}
	response = c.userService.ListUserReferral(request)
	ctx.JSON(response.HttpCode, response)
}

func deserializeGetUserReferralRequest(request interface{}) (userRequestDTO.ListUserReferralRequestDTO, error) {
	otpDTO := request.(userRequestDTO.ListUserReferralRequestDTO)

	cipheTextSalesEmail, err := base64.StdEncoding.DecodeString(otpDTO.SalesEmail)
	if err != nil {
		fmt.Println(err.Error())
	}
	plainTextSalesEmail, err := helper.RsaDecryptFromFEInBE([]byte(cipheTextSalesEmail))
	if err != nil {
		fmt.Println(err.Error())
		return userRequestDTO.ListUserReferralRequestDTO{}, err
	}

	var result userRequestDTO.ListUserReferralRequestDTO

	result.SalesEmail = plainTextSalesEmail
	result.Limit = otpDTO.Limit
	result.Offset = otpDTO.Offset
	result.Keyword = otpDTO.Keyword
	result.StartDate = otpDTO.StartDate
	result.EndDate = otpDTO.EndDate
	return result, nil
}
