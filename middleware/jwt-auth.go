package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-api/dto/response"
	"go-api/helper"
	"go-api/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	var response response.UniversalResponse
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		defer func() {
			if err := recover(); err != nil {
				response.HttpCode = http.StatusBadRequest
				response.ResponseCode = "99"
				response.ResponseMessage = "Unchaught error"
				response.Data = nil
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}()
		if authHeader == "" {
			// response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
			response.HttpCode = http.StatusBadRequest
			response.ResponseCode = "99"
			response.ResponseMessage = "No token found"
			response.Data = nil
			c.AbortWithStatusJSON(response.HttpCode, response)
			return
		}
		token, err := jwtService.ValidateToken(authHeader)
		if token == nil {
			log.Println(err)
			// response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
			response.HttpCode = http.StatusBadRequest
			response.ResponseCode = "99"
			response.ResponseMessage = "Token is not valid"
			response.Data = nil
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claim[user_id]: ", claims["user_id"])
			log.Println("Claim[issuer] :", claims["issuer"])
			log.Println("Claim[Expires_at]", claims["exp"])

			expiration := int64(claims["exp"].(float64))
			expirationTime := time.Unix(expiration, 0)
			expirationFormatted := expirationTime.Format("02 January 2006 15:04:05 MST")

			fmt.Println("Token kedaluwarsa pada:", expirationFormatted)

			if time.Now().Unix() > expiration {
				response.HttpCode = http.StatusBadRequest
				response.ResponseCode = "99"
				response.ResponseMessage = "Token is expired"
				response.Data = nil
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}

		}
	}
}

func AuthorizeJWTFinal(jwtService service.JWTService, email, bearerToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeaderRequest := c.GetHeader("Authorization")
		if authHeaderRequest == "" {
			response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
	}
}
