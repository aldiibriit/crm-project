package service

import (
	"log"

	"go-api/dto"
	"go-api/entity"
	internalRepository "go-api/repository/internal-repo"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is a contract about something that this service can do
type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	SaveToken(email, token string) error
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository internalRepository.UserRepository
	jwtService     JWTService
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRep internalRepository.UserRepository, jwtServ JWTService) AuthService {
	return &authService{
		userRepository: userRep,
		jwtService:     jwtServ,
	}
}

func (service *authService) VerifyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)
	var res2 entity.User
	res2 = res.(entity.User)
	if v, ok := res.(entity.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			res2.Name = v.Name
			res2.Email = v.Email
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
		log.Println("Failed map %v", err)
	}
	res := service.userRepository.InsertUser(userToCreate)
	return res
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func (service *authService) SaveToken(email, token string) error {
	return service.userRepository.SaveToken(email, token)
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
