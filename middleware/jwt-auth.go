package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"go-api/helper"
	"go-api/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		token, err := jwtService.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claim[user_id]: ", claims["user_id"])
			log.Println("Claim[issuer] :", claims["issuer"])
		} else {
			log.Println(err)
			response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}

func AuthorizeJWT2(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeaderRequest := c.GetHeader("Authorization")
		if authHeaderRequest == "" {
			response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		cipherText, _ := base64.StdEncoding.DecodeString(authHeaderRequest[7:])
		authHeaderDecrypted, _ := helper.RsaDecryptFromFEInBE([]byte(cipherText))
		token, err := jwtService.ValidateToken(authHeaderDecrypted)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claim[user_id]: ", claims["user_id"])
			log.Println("Claim[issuer] :", claims["issuer"])
			log.Println("Claim[expires_at] : ", claims["exp"])
			second := time.Duration(claims["exp"].(float64) * float64(time.Second))
			log.Println("expires in : ", second, " seconds")
			log.Println(claims)
		} else {
			log.Println(err)
			response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
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
