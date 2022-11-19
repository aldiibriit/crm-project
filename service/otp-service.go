package service

import (
	"fmt"
	"go-api/dto/request/emailRequestDTO"
	"go-api/dto/request/otpRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/entity"
	"golang_api-main/repository"
	"log"
	"reflect"
	"time"

	"github.com/google/uuid"
)

const otpRedaksi1 = "Nasabah Yth,\n" +
	"Kode Otorisasi untuk Platform Homespot anda adalah: "
const otpRedaksi2 = ". Kode anda berlaku 2 menit.\n" +
	"\n" +
	"Jangan Berikan Kode Rahasia ini ke Siapapun. Penyalahgunaan Kode Rahasia ini Menjadi Tanggung Jawab Pengguna"

type OTPService interface {
	SendOTP(email string) string
	ValidateOTP(request otpRequestDTO.ValidateOTPRequest) responseDTO.Response
	updateOTP(otp entity.TblOtp) bool
}

type otpService struct {
	otpRepository        repository.OTPRepository
	otpAttemptRepository repository.OTPAttemptRepository
	emailService         EmailService
}

func NewOTPService(otpRepo repository.OTPRepository, emailServ EmailService, otpAttemptRepo repository.OTPAttemptRepository) OTPService {
	return &otpService{
		otpRepository:        otpRepo,
		emailService:         emailServ,
		otpAttemptRepository: otpAttemptRepo,
	}
}

func (service *otpService) SendOTP(email string) string {
	start := time.Now()
	expired := start.Add(2 * time.Minute)
	generatedOtp := fmt.Sprint(time.Now().Nanosecond())
	otpCode := generatedOtp[:6]
	createdOTP := entity.TblOtp{
		ID:         uuid.New().String(),
		MobileNo:   "",
		Email:      email,
		OTP:        otpCode,
		OtpAttempt: 0,
		Status:     "ON_REQUESTED",
		CreatedAt:  start,
		UpdatedAt:  start,
		ExpiredAt:  expired,
	}

	service.otpRepository.UpdateOrCreate(createdOTP)
	log.Println("On sending OTP email")

	emailRequest := emailRequestDTO.EmailRequestDTO{
		Subject:   "Verifikasi OTP Homespot",
		Action:    4,
		EmailBody: otpRedaksi1 + otpCode + otpRedaksi2,
		ToAddres:  email,
	}

	service.emailService.SendMessage(emailRequest)

	return generatedOtp
}

func (service *otpService) ValidateOTP(request otpRequestDTO.ValidateOTPRequest) responseDTO.Response {
	var response responseDTO.Response
	otp := service.otpRepository.FindByOTPAndStatus(request, "ON_REQUESTED")
	if otp.Email == "" {
		log.Println("OTP Not Found")
		log.Println("Error Invalid OTP[Modul=OtpService|Method=ValidateOTP")

		otpAttempt := service.otpAttemptRepository.FindByEmail(request.Email)

		if otpAttempt.Email == "" {
			log.Println("Create new otp attempt")
			otpAttempt = entity.TblOtpAttempt{
				Id:         uuid.New().String(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				Email:      request.Email,
				OtpAttempt: 1,
			}
		} else {
			log.Println("Object OTP Attempt Found")
			currentTime := time.Now()
			otpAttempt.UpdatedAt = currentTime
			difference := currentTime.Sub(otpAttempt.UpdatedAt)
			dateDiff := int64(difference.Hours() / 24)
			if dateDiff > 3 {
				log.Println("Delete OTP Attempt Expired > 3 days")
				service.otpAttemptRepository.Delete(otpAttempt.Email)
				otpAttempt = entity.TblOtpAttempt{
					Id:         uuid.New().String(),
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Email:      request.Email,
					OtpAttempt: 1,
				}
			} else if dateDiff <= 3 {
				otpAttempt.OtpAttempt += 1
				if otpAttempt.OtpAttempt > 3 {
					log.Println("Otp Attempt Reach > 3x")
					response.HttpCode = 422
					response.MetadataResponse = nil
					response.ResponseCode = "99"
					response.ResponseDesc = "Reached limit for OTP input"
					response.ResponseData = nil
					response.Summary = nil
					return response
				}
			}
		}
		service.otpAttemptRepository.UpdateOrCreate(otpAttempt)
		log.Println("Success Update OTP Attempt")
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "OTP Not Valid"
		response.ResponseData = nil
		response.Summary = nil
		return response
	} else {
		log.Println("OTP Found")
		otpAttempt := service.otpAttemptRepository.FindByEmail(request.Email)
		log.Println("Check Null => ", reflect.ValueOf(otpAttempt).IsZero())
		if otpAttempt.Email != "" && otpAttempt.OtpAttempt >= 3 {
			log.Println("OTP Attempt FOUND")
			log.Println("Attempts : ", otpAttempt.OtpAttempt)
			response.HttpCode = 422
			response.MetadataResponse = nil
			response.ResponseCode = "99"
			response.ResponseDesc = "Reach limit for OTP input"
			response.ResponseData = nil
			response.Summary = nil
		}
	}

	log.Println("Clear Area")
	if !service.updateOTP(otp) {
		log.Println("More than 2 minutes")
		response.HttpCode = 422
		response.MetadataResponse = nil
		response.ResponseCode = "99"
		response.ResponseDesc = "OTP Expired"
		response.ResponseData = nil
		response.Summary = nil
		return response
	}

	log.Println("Less than 2 minutes")
	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = nil
	response.Summary = nil

	return response
}

func (service *otpService) updateOTP(otp entity.TblOtp) bool {
	log.Println("Start method updateOTP")
	var returnValue bool

	currentTime := time.Now()
	otp.UpdatedAt = currentTime
	log.Println("Updated At OTP : ", otp.UpdatedAt)
	log.Println("currentTime : ", currentTime)
	difference := currentTime.Sub(otp.CreatedAt)
	delta := difference.Minutes()

	log.Println("Validated : [Update OTP] => ", delta)
	if delta > 2 {
		otp.Status = "ON_EXPIRED"
		returnValue = false
	} else {
		otp.Status = "ON_VALIDATED"
		returnValue = true
	}
	otp.OtpAttempt += 1
	service.otpRepository.UpdateOrCreate(otp)

	return returnValue

}
