package service

import (
	"fmt"
	"go-api/dto"
	"go-api/dto/request/userRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/entity"
	"go-api/repository"

	"github.com/mashingan/smapping"
)

// UserService is a contract.....
type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userID string) entity.User
	GetDeveloper(request userRequestDTO.ListUserDeveloperRequestDTO) responseDTO.Response
}

type userService struct {
	userRepository repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		fmt.Println("Failed map %v:", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) entity.User {
	return service.userRepository.ProfileUser(userID)
}

func (service *userService) GetDeveloper(request userRequestDTO.ListUserDeveloperRequestDTO) responseDTO.Response {
	var response responseDTO.Response
	data := service.userRepository.GetDeveloper(request)
	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "Success"
	response.ResponseData = data
	response.Summary = nil
	return response
}
