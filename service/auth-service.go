package service

import (
	"encoding/base64"
	"log"

	"go-api/dto"
	"go-api/entity"
	"go-api/helper"
	"go-api/repository"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is a contract about something that this service can do
type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRep repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRep,
	}
}

func (service *authService) VerifyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)
	var res2 entity.User
	res2 = res.(entity.User)
	if v, ok := res.(entity.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			encryptedName, _ := helper.RsaEncryptBEToFE([]byte(v.Name))
			encryptedEmail, _ := helper.RsaEncryptBEToFE([]byte(v.Email))
			encodingName := base64.StdEncoding.EncodeToString([]byte(encryptedName))
			encodingEmail := base64.StdEncoding.EncodeToString([]byte(encryptedEmail))
			res2.Name = encodingName
			res2.Email = encodingEmail
			return res2
		}
		return false
	}
	return false
}

func (service *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.userRepository.InsertUser(userToCreate)
	return res
}

func (service *authService) FindByEmail(email string) entity.User {
	return service.userRepository.FindByEmail(email)
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
