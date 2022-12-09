package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-api/dto"
	"go-api/dto/request/authRequestDTO"
	"go-api/dto/request/emailRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/dto/response/authResponseDTO"
	"go-api/entity"
	"go-api/helper"
	"go-api/repository"

	"github.com/dgrijalva/jwt-go"
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
	PassthroughLogin(request authRequestDTO.PassthroughLoginRequest) (responseDTO.Response, string)
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository         repository.UserRepository
	salesRepository        repository.SalesRepository
	emailAttemptRepository repository.EmailAttemptRepository
	emailService           EmailService
	otpService             OTPService
	jwtService             JWTService
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRep repository.UserRepository, salesRepo repository.SalesRepository, emailServ EmailService, emailAttemptRepo repository.EmailAttemptRepository, otpServ OTPService, jwtServ JWTService) AuthService {
	return &authService{
		userRepository:         userRep,
		salesRepository:        salesRepo,
		emailService:           emailServ,
		emailAttemptRepository: emailAttemptRepo,
		otpService:             otpServ,
		jwtService:             jwtServ,
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
		fmt.Println("Failed map %v", err)
	}
	res := service.userRepository.InsertUser(userToCreate)
	return res
}

func (service *authService) CreateUserSales(request dto.RegisterSalesDTO) responseDTO.Response {
	var response responseDTO.Response

	LatestID := service.userRepository.GetLatestId()

	id := LatestID.ID + 1

	userToCreate := entity.TblUser{
		Email:          request.EmailSales,
		RegistrationId: uuid.New().String(),
		Type:           "Sales",
		Status:         "Registered",
		ID:             id,
		IdResponse:     strconv.Itoa(id),
		MobileNo:       request.SalesPhone,
		ModifiedAt:     time.Now(),
	}

	salesToCreate := entity.TblSales{
		EmailDeveloper: request.EmailDeveloper,
		EmailSales:     request.EmailSales,
		RegisteredBy:   request.RegisteredBy,
		SalesName:      request.SalesName,
		RefferalCode:   helper.GenerateRefferalCode(6),
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
	urlEncoded := "http://172.18.2.94/sales/activate/" + url.QueryEscape(urlEncrypted)
	emailRequest := emailRequestDTO.EmailRequestDTO{
		ToAddres:   request.EmailSales,
		UrlEncoded: urlEncoded,
		Action:     1,
	}

	if !service.emailService.SendMessage(emailRequest) {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "Cannot send message to " + emailRequest.ToAddres
		response.ResponseData = nil
		response.Summary = nil
		return response
	}

	encryptedCreatedUser := serializeCreatedUser(userToCreate)

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = encryptedCreatedUser
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

	intAction, err := strconv.Atoi(request.Action)
	if err != nil {
		log.Println("Error when convert action to string at auth service line 175 ", err.Error())
	}

	if intAction == 1 || intAction == 2 {
		emailAttempt := service.emailAttemptRepository.FindByEmailAndAction(request.Email, intAction)
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

	userIdentity := serializeActivatedUser(user)

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = userIdentity
	response.Summary = nil

	return response
}

func (service *authService) FindByEmail(email string) entity.User {
	return service.userRepository.FindByEmail(email)
}

func (service *authService) PassthroughLogin(request authRequestDTO.PassthroughLoginRequest) (responseDTO.Response, string) {
	var response responseDTO.Response
	var responseCreateTokenMap map[string]interface{}
	var responseLoginMap map[string]interface{}

	urlCreateToken := "https://api.homespot.id/mes/api/v1/auth/createToken"
	urlLogin := "https://api.homespot.id/mes/api/v1/user/login"
	method := "POST"

	payloadCreateToken := strings.NewReader(`{` + "" + `"email" : "` + request.Email + `",` + "" + `
    "applicationName" : "HOMESPOT"` + "" + `}`)

	payloadLogin := strings.NewReader(`{` + "" + `"email" : "` + request.Email + `",` + "" + `"password" : "` + request.Password + `"` + "" + `}`)

	responseAPICreateToken := callAPICreateToken(method, urlCreateToken, payloadCreateToken)
	json.Unmarshal([]byte(responseAPICreateToken), &responseCreateTokenMap)
	if responseCreateTokenMap["responseCode"] != "00" {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = "Error in call API Create Token"
		response.Summary = nil
		return response, ""
	}

	encryptedTokenFromBE := responseCreateTokenMap["responseData"].([]interface{})[0]
	strEncryptToken := fmt.Sprintf("%v", encryptedTokenFromBE)

	decodedToken, err := base64.StdEncoding.DecodeString(strEncryptToken)
	if err != nil {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = err.Error()
		response.Summary = nil
		return response, ""
	}
	decryptedToken, err := helper.RsaDecryptFromBEInFE(decodedToken)
	if err != nil {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = err.Error()
		response.Summary = nil
		return response, ""
	}

	fixToken, err := helper.RsaEncryptFEToBE([]byte(decryptedToken))
	if err != nil {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = err.Error()
		response.Summary = nil
		return response, ""
	}

	expiredAt := extractToken(decryptedToken)

	bearerToken := "Bearer " + fixToken

	responseLogin := CallAPILogin(method, urlLogin, payloadLogin, bearerToken)
	json.Unmarshal([]byte(responseLogin), &responseLoginMap)
	if responseLoginMap["responseCode"] != "00" {
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseData = nil
		response.ResponseDesc = "Error in call API Login"
		response.Summary = nil
		return response, ""
	}

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = nil
	response.ResponseDesc = "Success"
	response.Summary = nil
	setCookieValue := `jwt=` + decryptedToken + `;Path=/;Expires=` + expiredAt + `;HttpOnly;Samesite=Lax;Domain=.homespot.id`

	return response, setCookieValue
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
	response.ResponseData = nil
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

func serializeCreatedUser(request interface{}) entity.TblUser {
	data := request.(entity.TblUser)

	encryptedIdResponse, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(data.ID)))
	encryptedEmail, _ := helper.RsaEncryptBEToFE([]byte(data.Email))
	encryptedMobileNo, _ := helper.RsaEncryptBEToFE([]byte(data.MobileNo))
	encryptedRegistrationId, _ := helper.RsaEncryptBEToFE([]byte(data.RegistrationId))
	encryptedStatus, _ := helper.RsaEncryptBEToFE([]byte(data.Status))
	encryptedType, _ := helper.RsaEncryptBEToFE([]byte(data.Type))
	encryptedCreatedAt, _ := helper.RsaEncryptBEToFE([]byte(data.CreatedAt.String()))
	encryptedModifiedAt, _ := helper.RsaEncryptBEToFE([]byte(data.ModifiedAt.String()))

	var result entity.TblUser

	result.IdResponse = encryptedIdResponse
	result.Email = encryptedEmail
	result.MobileNo = encryptedMobileNo
	result.RegistrationId = encryptedRegistrationId
	result.Status = encryptedStatus
	result.Type = encryptedType
	result.CreatedAtRes = encryptedCreatedAt
	result.ModifiedAtRes = encryptedModifiedAt
	return result
}

func serializeActivatedUser(request interface{}) authResponseDTO.UserIndetityDtoRes {
	data := request.(entity.TblUser)

	encryptedUsername, _ := helper.RsaEncryptBEToFE([]byte(data.Email))

	var result authResponseDTO.UserIndetityDtoRes

	result.Username = encryptedUsername
	return result
}

func callAPICreateToken(method string, url string, payload *strings.Reader) string {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err.Error()
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	return string(body)
}

func CallAPILogin(method string, url string, payload *strings.Reader, bearerToken string) string {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(body)
}

func extractToken(tokenString string) string {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Can't convert token's claims to standard claims")
	}

	var tm time.Time
	switch iat := claims["iat"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}

	return tm.Format(time.RFC1123)
}
