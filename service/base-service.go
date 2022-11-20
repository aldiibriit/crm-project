package service

import (
	"encoding/base64"
	"go-api/entity"
	"go-api/helper"
	"go-api/repository"
	"log"
)

type BaseService interface {
	ValidateToken(bearerToken, email string) (bool, string)
}

type baseService struct {
	jwtHistRepository repository.JwtHistRepository
	jwtService        JWTService
}

func NewBaseService(jwtHistRepo repository.JwtHistRepository, jwtServ JWTService) BaseService {
	return &baseService{
		jwtHistRepository: jwtHistRepo,
		jwtService:        jwtServ,
	}
}

func (service *baseService) ValidateToken(bearerToken, email string) (bool, string) {
	token := ""
	bearerTokenStartWith := bearerToken

	// decrypt from FE
	cipherText, _ := base64.StdEncoding.DecodeString(bearerTokenStartWith[7:])
	decryptedBearerToken, _ := helper.RsaDecryptFromFEInBE(cipherText)

	// start validate is exist?
	token = decryptedBearerToken
	jwtHist := service.jwtHistRepository.FindByEmail(email)

	var emptyJwtHist entity.JwtHistGo

	// end validate is exist?
	if jwtHist == emptyJwtHist {
		log.Println("Error JwtHist Not Found !")
		return false, "Error JwtHist Not Found !"
	}

	// validate is matching?
	if jwtHist.Jwt != token {
		log.Println("Error JwtHist Doesn't Match")
		return false, "Error JwtHist Doesn't Match"
	}

	jwtToken, err := service.jwtService.ValidateToken(jwtHist.Jwt)
	if jwtToken.Valid {
		return true, ""
	} else {
		log.Println(err.Error())
		return false, err.Error()
	}

	return false, err.Error()
}
