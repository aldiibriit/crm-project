package service

import (
	"fmt"
	"go-api/dto/request/emailRequestDTO"
	"go-api/entity"
	"go-api/repository"
	"log"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/google/uuid"
)

type EmailService interface {
	SendMessage(request emailRequestDTO.EmailRequestDTO) bool
	sendMail(request emailRequestDTO.EmailRequestDTO) bool
	limitExceed(email string, action int, urlEncoded string) bool
}

type emailService struct {
	emailAttemptRepository repository.EmailAttemptRepository
}

func NewEmailService(emailAttemptRepo repository.EmailAttemptRepository) EmailService {
	return &emailService{
		emailAttemptRepository: emailAttemptRepo,
	}
}

func (service *emailService) SendMessage(request emailRequestDTO.EmailRequestDTO) bool {
	return service.sendMail(request)
}

func (service *emailService) sendMail(request emailRequestDTO.EmailRequestDTO) bool {

	if service.limitExceed(request.ToAddres, request.Action, request.UrlEncoded) {
		log.Println("Error in limit exceed")
		return false
	} else {
		dailer := mail.NewDialer("email-smtp.ap-southeast-1.amazonaws.com", 587, "AKIA4RGQ2AWFS7BGLKOX", "BHa+AUiY9lqKitS9qAzregJsPwynZswwoXL4a314szp1")
		dailer.Timeout = 5 * time.Second

		msg := mail.NewMessage()
		msg.SetHeader("To", request.ToAddres)
		msg.SetHeader("Subject", request.Subject)
		msg.SetBody("text/html", generateBody(request))
		msg.SetHeader("From", "noreply2@homespot.id")

		err := dailer.DialAndSend(msg)
		if err != nil {
			log.Println(err.Error())
		}

		return true
	}
}

func (service *emailService) limitExceed(email string, action int, urlEncoded string) bool {
	log.Println("Start M=limitExceed S=EmailService")
	nowTs := time.Now()

	emailAttempt := service.emailAttemptRepository.FindByEmailAndAction(email, action)
	if emailAttempt.Id != "" {
		log.Println("Email attempt not null")
		difference := nowTs.Sub(emailAttempt.CreatedAt)
		delta := int64(difference.Hours() / 24)
		fmt.Println("Hour : ", delta)
		if delta > 0 {
			log.Println("Reset to zero")
			emailAttempt.Attempt = 0
		} else {
			emailAttempt.Attempt += 1
			if emailAttempt.Attempt > 3 {
				log.Println("Attempt 1 : ", emailAttempt.Attempt)
				return true
			}
			log.Println("Attempt 2 : ", emailAttempt.Attempt)
		}
		emailAttempt.UrlEncoded = urlEncoded
		emailAttempt.UpdatedAt = nowTs
		// service.emailAttemptRepository.UpdateOrCreate(emailAttempt)
		// return false
	} else {
		log.Println("Email Attempt is null")
		emailAttempt = entity.TblEmailAttempt{
			Attempt:    1,
			Id:         uuid.New().String(),
			Email:      email,
			UrlEncoded: urlEncoded,
			Action:     action,
			CreatedAt:  nowTs,
			UpdatedAt:  nowTs,
		}
		// log.Println("Data : ", emailAttempt)
	}

	service.emailAttemptRepository.UpdateOrCreate(emailAttempt)
	log.Println("End M=Check limit exceed")
	return false
}

func generateBody(emailRequest emailRequestDTO.EmailRequestDTO) string {
	body := ""
	if emailRequest.Action == 3 {
		body = "<html><head><style>.doneRegis_title{font-weight:700;font-size:2rem}.doneRegisifAsk{text-align:center;margin-top:2rem}.doneRegisconWrap{padding:.2rem;background:#ebebeb;text-align:center;align-items:center}.doneRegistblWrap{margin:0 auto;text-align:center}.doneRegistblTxt{margin-top:0;margin-bottom:0;font-weight:400}</style></head><body><div><h1 class=\"doneRegistitle\">Registrasi Kamu Berhasil</h1><p>Selamat<b>" + emailRequest.Name + "</b>, Kamu bisa mencari dan membeli rumah idamanmu melalui Homespot</p><p class=\"doneRegisifAsk\">Bila ada pertanyaan, Silahkan menghubungi kami :</p><div class=\"doneRegisconWrap\"><table class=\"doneRegistblWrap\"><tr><th><img src=\"https://storage.googleapis.com/artifacts.concrete-plasma-244309.appspot.com/homespot/wa-tiny.png\" alt=\"wa-tiny-icon\"></th><th><p class=\"doneRegistblTxt\">+622150864230</p></th></tr></table><table class=\"doneRegistblWrap\"><tr><th><img src=\"https://storage.googleapis.com/artifacts.concrete-plasma-244309.appspot.com/homespot/email-tiny.png\" alt=\"email-tiny-icon\"></th><th><p class=\"doneRegis_tblTxt\">support@homespot.co.id</p></th></tr></table></div></div></body></html>"
	} else if emailRequest.Action == 4 || emailRequest.Action == 5 { //otp activation
		return emailRequest.EmailBody
	} else {
		body = "<html><head><link href=\"https://cdn.jsdelivr.net/npm/remixicon@2.5.0/fonts/remixicon.css\" rel=\"stylesheet\"><style>.logo{margin-bottom:.7rem}.verify_wrap{text-align:center}.verifytitle{font-weight:700;font-size:2rem}.verifylink{font-weight:700;text-decoration:underline}</style></head><body><div><div class=\"verifywrap\"><img class=\"logo\" src=\"https://storage.googleapis.com/artifacts.concrete-plasma-244309.appspot.com/homespot/logo/logo.svg\"><h1 class=\"verifytitle\">Verifikasi Email Kamu</h1><p>Data registrasi kamu telah berhasil kami terima. Verifikasi email kamu dengan mengklik tautan di bawah ini:</p><p class=\"verifylink\"><a href=\"" + emailRequest.UrlEncoded + "\">Verifikasi Email</a></p><p>Atau kamu dapat menyalin link berikut untuk memverifikasi email kamu</p><p class=\"verify_link\"><a href=\"" + emailRequest.UrlEncoded + "\">" + emailRequest.UrlEncoded + "</a></p></div></div></body></html>"
	}
	return body
}
