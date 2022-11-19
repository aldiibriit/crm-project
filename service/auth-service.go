package service

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"time"

	"go-api/dto"
	"go-api/dto/request/authRequestDTO"
	"go-api/dto/request/emailRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/entity"
	"go-api/helper"
	"go-api/repository"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is a contract about something that this service can do
type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	CreateUserSales(user dto.RegisterSalesDTO) responseDTO.Response
	ActivateUser(request authRequestDTO.ActivateRequestDTO) responseDTO.Response
	PasswordConfirmation(request dto.PasswordConfirmationDTO) responseDTO.Response
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository         repository.UserRepository
	salesRepository        repository.SalesRepository
	emailAttemptRepository repository.EmailAttemptRepository
	emailService           EmailService
	otpService             OTPService
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRep repository.UserRepository, salesRepo repository.SalesRepository, emailServ EmailService, emailAttemptRepo repository.EmailAttemptRepository, otpServ OTPService) AuthService {
	return &authService{
		userRepository:         userRep,
		salesRepository:        salesRepo,
		emailService:           emailServ,
		emailAttemptRepository: emailAttemptRepo,
		otpService:             otpServ,
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

func (service *authService) CreateUserSales(request dto.RegisterSalesDTO) responseDTO.Response {
	var response responseDTO.Response

	userToCreate := entity.TblUser{
		Email:          request.EmailSales,
		RegistrationId: uuid.New().String(),
		Type:           "Sales",
		Status:         "Registered",
		ID:             666,
		MobileNo:       request.SalesPhone,
		ModifiedAt:     time.Now(),
	}

	salesToCreate := entity.TblSales{
		EmailDeveloper: request.EmailDeveloper,
		EmailSales:     request.EmailSales,
		RegisteredBy:   request.RegisteredBy,
		RefferalCode:   uuid.New().String(),
		ModifiedAt:     time.Now(),
	}

	err := service.userRepository.InsertUserSales(userToCreate)
	if err != nil {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "Failed insert to tbl_user"
		response.ResponseData = nil
		response.Summary = nil
		return response
	}
	err = service.salesRepository.InsertRelation(salesToCreate)
	if err != nil {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "Failed insert to tbl_user_sales"
		response.ResponseData = nil
		response.Summary = nil
		return response
	}

	urlEncrypted, _ := helper.RsaEncryptFEToBE([]byte(request.EmailSales))
	urlEncoded := "http://dev.homespot.id/user/activate" + url.QueryEscape(urlEncrypted)
	emailRequest := emailRequestDTO.EmailRequestDTO{
		ToAddres:   request.EmailSales,
		UrlEncoded: urlEncoded,
		Action:     1,
	}

	fmt.Println(emailRequest)

	if !service.emailService.SendMessage(emailRequest) {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "Cannot send message to " + emailRequest.ToAddres
		response.ResponseData = nil
		response.Summary = nil
		return response
	}

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = salesToCreate
	response.Summary = nil

	return response
}

func (service *authService) ActivateUser(request authRequestDTO.ActivateRequestDTO) responseDTO.Response {
	log.Println("Request : ", request)
	var response responseDTO.Response
	fmt.Println("Start[Modul=AuthService|Method=activateUser|Data=", request.Action)
	fmt.Println("Validate user by email,D=[Email=", request.Email, "|Action=", request.Action)
	user := service.userRepository.FindByEmail2(request.Email)
	if user.Email == "" {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "User not found !"
		response.ResponseData = nil
		response.Summary = nil
		return response
	}

	if request.Action == 1 || request.Action == 2 {
		emailAttempt := service.emailAttemptRepository.FindByEmailAndAction(request.Email, request.Action)
		fmt.Println(emailAttempt)
		if emailAttempt.Email == "" {
			response.HttpCode = 422
			response.MetadataResponse = nil
			response.ResponseCode = "99"
			response.ResponseDesc = "Email attempt not found"
			response.ResponseData = nil
			response.Summary = nil
			return response
		}

		log.Println("Start A=Validate URL ENCODED")

		if request.UrlEncoded != emailAttempt.UrlEncoded {
			response.HttpCode = 422
			response.MetadataResponse = nil
			response.ResponseCode = "99"
			response.ResponseDesc = "Url encoded not valid"
			response.ResponseData = nil
			response.Summary = nil
			return response
		}

		now := time.Now()
		lastModified := emailAttempt.UpdatedAt
		difference := now.Sub(lastModified)
		delta := int64(difference.Hours() / 24)

		log.Println("Message[Date Diff]=", delta)
		if delta > 3 {
			response.HttpCode = 422
			response.MetadataResponse = nil
			response.ResponseCode = "99"
			response.ResponseDesc = "Email expired"
			response.ResponseData = nil
			response.Summary = nil
			return response
		}

		if user.Status != "Registered" {
			response.HttpCode = 422
			response.MetadataResponse = nil
			response.ResponseCode = "99"
			response.ResponseDesc = "Data not found"
			response.ResponseData = nil
			response.Summary = nil
			return response
		}

		otp := service.otpService.SendOTP(request.Email)

		if otp == "" {
			response.HttpCode = 422
			response.MetadataResponse = nil
			response.ResponseCode = "99"
			response.ResponseDesc = "General Error (Send OTP)"
			response.ResponseData = nil
			response.Summary = nil
		}

		log.Println("OTP Sent")

	} else {
		log.Println("Update User ACTIVE")
		user.CreatedAt = time.Now()
		user.ModifiedAt = time.Now()
		user.Status = "Active"

		service.userRepository.UpdateOrCreate(user)
	}

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = user
	response.Summary = nil

	return response
}

func (service *authService) FindByEmail(email string) entity.User {
	return service.userRepository.FindByEmail(email)
}

func (service *authService) PasswordConfirmation(request dto.PasswordConfirmationDTO) responseDTO.Response {
	var response responseDTO.Response
	user := service.userRepository.FindByEmail2(request.Email)
	if user.Email == "" {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "User not found"
		response.ResponseData = nil
		response.Summary = nil
		return response
	}
	user.Password = request.NewPassword
	user.ModifiedAt = time.Now()
	updatedUser := entity.TblUser{
		Email:    request.Email,
		Password: request.NewPassword,
	}
	service.userRepository.UpdateOrCreate(updatedUser)
	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = updatedUser
	response.Summary = nil
	return response
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
