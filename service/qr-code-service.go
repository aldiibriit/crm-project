package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	QrCodeRequest "go-api/dto/request/qr-code"
	"go-api/dto/response"
	ExternalRepository "go-api/repository/external-repo"
	InternalRepository "go-api/repository/internal-repo"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
)

type QrCodeService interface {
	GenerateQr(request QrCodeRequest.GenerateQr) response.UniversalResponse
}

type qrCodeService struct {
	minioRepository       ExternalRepository.MinioRepository
	logActivityRepository InternalRepository.LogActivityRepository
}

func NewQrCodeService(minioRepo ExternalRepository.MinioRepository, logActivityRepo InternalRepository.LogActivityRepository) QrCodeService {
	return &qrCodeService{
		minioRepository:       minioRepo,
		logActivityRepository: logActivityRepo,
	}
}

func (service *qrCodeService) GenerateQr(request QrCodeRequest.GenerateQr) response.UniversalResponse {
	var response response.UniversalResponse

	// get data to encoded
	jsonByte, err := json.Marshal(request)
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ERROR_MARSHAL_MESSAGE")
		response.Data = nil
		return response
	}

	//generate qr
	var png []byte
	png, err = qrcode.Encode(string(jsonByte), qrcode.Medium, 256)
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ENCODE_ERROR")
		response.Data = nil
		return response
	}

	reader := bytes.NewReader(png)

	err = service.minioRepository.UploadFile(reader, reader.Size(), "crm-project", "qr-code/"+request.Sn+".png")
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("MINIO_FAILED_UPLOAD_FILE_MESSAGE")
		response.Data = nil
		return response
	}

	objectURL := fmt.Sprintf("%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), "crm-project", "qr-code/"+request.Sn+".png")
	successMap := map[string]interface{}{}
	successMap["url"] = objectURL

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = successMap

	return response
}
