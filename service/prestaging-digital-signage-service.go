package service

import (
	"encoding/json"
	logActivityRequest "go-api/dto/request/log-activity"
	prestagingDigitalSignageRequest "go-api/dto/request/prestaging-digital-signage"
	"go-api/dto/response"
	"go-api/entity"
	ExternalRepository "go-api/repository/external-repo"
	InternalRepository "go-api/repository/internal-repo"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/mashingan/smapping"
)

type PrestagingDigitalSignageService interface {
	PostPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingDigitalSignageRequest.PostPrestaging) response.UniversalResponse
	ApprovePrestaging(request prestagingDigitalSignageRequest.ApprovePrestaging) response.UniversalResponse
	RejectPrestaging(request prestagingDigitalSignageRequest.RejectPrestaging) response.UniversalResponse
	ReuploadPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingDigitalSignageRequest.PostPrestaging) response.UniversalResponse
	AllSubmittedData() response.UniversalResponse
	GetSubmittedDataBySn(request prestagingDigitalSignageRequest.FindBySn) response.UniversalResponse
	GetRejectedData(request prestagingDigitalSignageRequest.FindRejectedData) response.UniversalResponse
}

type prestagingDigitalSignageService struct {
	minioRepository                    ExternalRepository.MinioRepository
	logActivityRepository              InternalRepository.LogActivityRepository
	prestagingDigitalSignageRepository InternalRepository.PrestagingDigitalSignageRepository
	baseRepository                     InternalRepository.BaseRepository
	stagingDigitalSignageRepository    InternalRepository.StagingDigitalSignageRepository
}

func NewPrestagingDigitalSignageService(minioRepo ExternalRepository.MinioRepository, logActivityRepo InternalRepository.LogActivityRepository, prestagingDigitalSignageRepo InternalRepository.PrestagingDigitalSignageRepository, baseRepo InternalRepository.BaseRepository, stagingDigitalSignageRepo InternalRepository.StagingDigitalSignageRepository) PrestagingDigitalSignageService {
	return &prestagingDigitalSignageService{
		minioRepository:                    minioRepo,
		baseRepository:                     baseRepo,
		logActivityRepository:              logActivityRepo,
		prestagingDigitalSignageRepository: prestagingDigitalSignageRepo,
		stagingDigitalSignageRepository:    stagingDigitalSignageRepo,
	}
}

func (service *prestagingDigitalSignageService) PostPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingDigitalSignageRequest.PostPrestaging) response.UniversalResponse {
	var response response.UniversalResponse
	var err error
	ch := make(chan interface{}, 1)
	go func() {
		// time.Sleep purposed to held the upload file process (expected : the upload file function will execute when response sent to front end,so the client doesn't need to wait for file upload)
		time.Sleep(5 * time.Second)
		var dataChanel = <-ch
		if dataChanel == nil && err == nil {
			// count the process
			now := time.Now()
			defer timeTrack(now, "upload file to minio")

			for _, v := range requestMap {
				openedFile, err := v.Open()
				if err != nil {
					log.Println(err.Error())
				}
				defer openedFile.Close()

				fileReader := io.Reader(openedFile)
				err = service.minioRepository.UploadFile(fileReader, v.Size, "crm-project", "prestaging-digital-signage/"+v.Filename)
				if err != nil {
					log.Println(err.Error())
				}
			}
		} else if dataChanel != nil {
			log.Println(dataChanel)
		} else if err != nil {
			log.Println(err.Error())
		}
	}()

	// panic handler
	defer func() {
		if r := recover(); r != nil {
			ch <- r
		} else {
			ch <- nil
		}
	}()

	// start transaction
	tx := service.baseRepository.StartTransaction()

	// mapping (fill) from request to entity
	tbPrestagingDigitalSignage := entity.TbPrestagingDigitalSignage{}
	err = smapping.FillStruct(&tbPrestagingDigitalSignage, smapping.MapFields(&request))
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("MAPPING_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	// insert to tb_prestaging_DigitalSignage
	err = service.prestagingDigitalSignageRepository.InsertWithTx(tbPrestagingDigitalSignage, tx)
	if err != nil {
		log.Println(err.Error())
		service.baseRepository.RollbackTransaction(tx)
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATABASE_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	// insert to log activity
	logActivityRequest := logActivityRequest.InsertRequest{
		Sn:                request.Sn,
		StatusDescription: os.Getenv("LOG_SUBMITTED"),
		Category:          os.Getenv("CATEGORY_PRESTAGING_DIGITAL_SIGNAGE"),
	}

	tbLogActivity := entity.TbLogActivity{}
	err = smapping.FillStruct(&tbLogActivity, smapping.MapFields(&logActivityRequest))
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("MAPPING_ERROR_MESSAGE")
		response.Data = nil
		return response
	}
	err = service.logActivityRepository.InsertWithTx(tbLogActivity, tx)
	if err != nil {
		log.Println(err.Error())
		service.baseRepository.RollbackTransaction(tx)
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATABASE_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	// commit the transaction if there isn't error
	service.baseRepository.CommitTransaction(tx)

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = request

	return response
}

func (service *prestagingDigitalSignageService) ApprovePrestaging(request prestagingDigitalSignageRequest.ApprovePrestaging) response.UniversalResponse {
	var response response.UniversalResponse

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	// insert to tbl log
	tx := service.baseRepository.StartTransaction()
	tblLogActivity := entity.TbLogActivity{
		Category:          os.Getenv("CATEGORY_PRESTAGING_DIGITAL_SIGNAGE"),
		Sn:                request.Sn,
		StatusDescription: os.Getenv("LOG_APPROVED"),
	}

	dataExist := service.stagingDigitalSignageRepository.FindBySn(request.Sn)
	if dataExist.Sn == "" || len(dataExist.Sn) == 0 {
		tbStagingDigitalSignage := entity.TbStagingDigitalSignage{
			Sn: request.Sn,
		}
		err := service.stagingDigitalSignageRepository.InsertWithTx(tbStagingDigitalSignage, tx)
		if err != nil {
			log.Println(err.Error())
			service.baseRepository.RollbackTransaction(tx)
			response.HttpCode = 500
			response.ResponseCode = "99"
			response.ResponseMessage = os.Getenv("DATABASE_ERROR_MESSAGE")
			response.Data = nil
			return response
		}

	}
	err := service.logActivityRepository.InsertWithTx(tblLogActivity, tx)
	if err != nil {
		log.Println(err.Error())
		service.baseRepository.RollbackTransaction(tx)
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATABASE_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	service.baseRepository.CommitTransaction(tx)

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = nil

	return response
}

func (service *prestagingDigitalSignageService) RejectPrestaging(request prestagingDigitalSignageRequest.RejectPrestaging) response.UniversalResponse {
	var response response.UniversalResponse

	// mapping from request to entity
	var updatedData map[string]interface{}
	data, err := json.Marshal(request)
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ERROR_MARSHAL_MESSAGE")
		response.Data = nil
		return response
	}
	err = json.Unmarshal(data, &updatedData)
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ERROR_UNMARSHAL_MESSAGE")
		response.Data = nil
		return response
	}

	snakeCaseMap := make(map[string]interface{})

	for key, v := range updatedData {
		snakeCaseMap[strcase.ToSnake(key)] = v
	}

	tx := service.baseRepository.StartTransaction()

	// update data
	err = service.prestagingDigitalSignageRepository.UpdateWithTx(snakeCaseMap, tx)
	if err != nil {
		log.Println(err.Error())
		service.baseRepository.RollbackTransaction(tx)
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATABASE_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	tbLogActivity := entity.TbLogActivity{
		Category:          os.Getenv("CATEGORY_PRESTAGING_DIGITAL_SIGNAGE"),
		Sn:                request.Sn,
		StatusDescription: os.Getenv("LOG_REJECTED"),
	}

	err = service.logActivityRepository.InsertWithTx(tbLogActivity, tx)
	if err != nil {
		log.Println(err.Error())
		service.baseRepository.RollbackTransaction(tx)
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATABASE_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	// commit the transaction
	service.baseRepository.CommitTransaction(tx)

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = nil

	return response
}

func (service *prestagingDigitalSignageService) ReuploadPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingDigitalSignageRequest.PostPrestaging) response.UniversalResponse {
	var response response.UniversalResponse
	var err error
	ch := make(chan interface{}, 1)
	go func() {
		// time.Sleep purposed to held the upload file process (expected : the upload file function will execute when response sent to front end,so the client doesn't need to wait for file upload)
		time.Sleep(5 * time.Second)
		var dataChanel = <-ch
		if dataChanel == nil && err == nil {
			// count the process
			now := time.Now()
			defer timeTrack(now, "upload file to minio")

			for _, v := range requestMap {
				openedFile, err := v.Open()
				if err != nil {
					log.Println(err.Error())
				}
				defer openedFile.Close()

				fileReader := io.Reader(openedFile)
				err = service.minioRepository.UploadFile(fileReader, v.Size, "crm-project", "prestaging-digital-signage/"+v.Filename)
				if err != nil {
					log.Println(err.Error())
				}
			}
		} else if dataChanel != nil {
			log.Println(dataChanel)
		} else if err != nil {
			log.Println(err.Error())
		}
	}()

	// panic handler
	defer func() {
		if r := recover(); r != nil {
			ch <- r
		} else {
			ch <- nil
		}
	}()

	// start transaction
	tx := service.baseRepository.StartTransaction()

	var updatedData map[string]interface{}

	data, err := json.Marshal(request)
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ERROR_MARSHAL_MESSAGE")
		response.Data = nil
		return response
	}

	err = json.Unmarshal(data, &updatedData)
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ERROR_UNMARSHAL_MESSAGE")
		response.Data = nil
		return response
	}

	camelMap := make(map[string]interface{})

	for key, v := range updatedData {
		camelMap[strcase.ToSnake(key)] = v
	}

	// update to tb_prestaging_DigitalSignage
	err = service.prestagingDigitalSignageRepository.UpdateWithTx(camelMap, tx)
	if err != nil {
		log.Println(err.Error())
		service.baseRepository.RollbackTransaction(tx)
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATABASE_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	// insert to log activity
	logActivityRequest := logActivityRequest.InsertRequest{
		Sn:                request.Sn,
		StatusDescription: os.Getenv("LOG_REUPLOADED"),
		Category:          os.Getenv("CATEGORY_PRESTAGING_DIGITAL_SIGNAGE"),
	}

	tbLogActivity := entity.TbLogActivity{}
	err = smapping.FillStruct(&tbLogActivity, smapping.MapFields(&logActivityRequest))
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("MAPPING_ERROR_MESSAGE")
		response.Data = nil
		return response
	}
	err = service.logActivityRepository.InsertWithTx(tbLogActivity, tx)
	if err != nil {
		log.Println(err.Error())
		service.baseRepository.RollbackTransaction(tx)
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATABASE_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	// commit the transaction if there isn't error
	service.baseRepository.CommitTransaction(tx)

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = request

	return response
}

func (service *prestagingDigitalSignageService) AllSubmittedData() response.UniversalResponse {
	var response response.UniversalResponse
	data := service.prestagingDigitalSignageRepository.FindAllSubmittedData()
	if len(data) == 0 {
		response.HttpCode = 404
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ERROR_NOT_FOUND_DATA_MESSAGE")
		response.Data = nil
		return response
	}

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = data
	return response
}

func (service *prestagingDigitalSignageService) GetSubmittedDataBySn(request prestagingDigitalSignageRequest.FindBySn) response.UniversalResponse {
	var response response.UniversalResponse
	data := service.prestagingDigitalSignageRepository.FindBySn(request.Sn)
	if len(data.Sn) == 0 {
		response.HttpCode = 404
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ERROR_NOT_FOUND_DATA_MESSAGE")
		response.Data = nil
		return response
	}

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = data
	return response
}

func (service *prestagingDigitalSignageService) GetRejectedData(request prestagingDigitalSignageRequest.FindRejectedData) response.UniversalResponse {
	var response response.UniversalResponse
	data := service.prestagingDigitalSignageRepository.FindRejectedData(request.IdUploader)
	if len(data) == 0 {
		response.HttpCode = 404
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("ERROR_NOT_FOUND_DATA_MESSAGE")
		response.Data = nil
		return response
	}

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = data
	return response
}
