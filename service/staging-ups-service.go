package service

import (
	"encoding/json"
	logActivityRequest "go-api/dto/request/log-activity"
	stagingUpsRequest "go-api/dto/request/staging-ups"
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

type StagingUPSService interface {
	PostStaging(requestMap map[string]*multipart.FileHeader, request stagingUpsRequest.PostStaging) response.UniversalResponse
	ApproveStaging(request stagingUpsRequest.ApproveStaging) response.UniversalResponse
	RejectStaging(request stagingUpsRequest.RejectStaging) response.UniversalResponse
	ReuploadStaging(requestMap map[string]*multipart.FileHeader, request stagingUpsRequest.PostStaging) response.UniversalResponse
	AllSubmittedData() response.UniversalResponse
	GetSubmittedDataBySn(request stagingUpsRequest.FindBySn) response.UniversalResponse
	GetRejectedData(request stagingUpsRequest.FindRejectedData) response.UniversalResponse
}

type stagingUpsService struct {
	minioRepository       ExternalRepository.MinioRepository
	logActivityRepository InternalRepository.LogActivityRepository
	baseRepository        InternalRepository.BaseRepository
	stagingUpsRepository  InternalRepository.StagingUPSRepository
}

func NewStagingUPSService(minioRepo ExternalRepository.MinioRepository, logActivityRepo InternalRepository.LogActivityRepository, baseRepo InternalRepository.BaseRepository, stagingUPSRepo InternalRepository.StagingUPSRepository) StagingUPSService {
	return &stagingUpsService{
		minioRepository:       minioRepo,
		baseRepository:        baseRepo,
		logActivityRepository: logActivityRepo,
		stagingUpsRepository:  stagingUPSRepo,
	}
}

func (service *stagingUpsService) PostStaging(requestMap map[string]*multipart.FileHeader, request stagingUpsRequest.PostStaging) response.UniversalResponse {
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
				err = service.minioRepository.UploadFile(fileReader, v.Size, "crm-project", "staging-ups/"+v.Filename)
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

	snakeCaseMap := make(map[string]interface{})

	for key, v := range updatedData {
		snakeCaseMap[strcase.ToSnake(key)] = v
	}

	// insert to tb_staging_ups
	dataExist := service.stagingUpsRepository.FindBySn(request.Sn)
	if dataExist.Sn == "" || len(dataExist.Sn) == 0 {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATA_NOT_FOUND_MESSAGE")
		response.Data = nil
		return response
	}

	err = service.stagingUpsRepository.UpdateWithTx(snakeCaseMap, tx)
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
		Category:          os.Getenv("CATEGORY_STAGING_UPS"),
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

func (service *stagingUpsService) ApproveStaging(request stagingUpsRequest.ApproveStaging) response.UniversalResponse {
	var response response.UniversalResponse

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	// insert to tbl log
	tx := service.baseRepository.StartTransaction()
	tblLogActivity := entity.TbLogActivity{
		Category:          os.Getenv("CATEGORY_STAGING_UPS"),
		Sn:                request.Sn,
		StatusDescription: os.Getenv("LOG_APPROVED"),
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

func (service *stagingUpsService) RejectStaging(request stagingUpsRequest.RejectStaging) response.UniversalResponse {
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
	err = service.stagingUpsRepository.UpdateWithTx(snakeCaseMap, tx)
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
		Category:          os.Getenv("CATEGORY_STAGING_UPS"),
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

func (service *stagingUpsService) ReuploadStaging(requestMap map[string]*multipart.FileHeader, request stagingUpsRequest.PostStaging) response.UniversalResponse {
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
				err = service.minioRepository.UploadFile(fileReader, v.Size, "crm-project", "staging-ups/"+v.Filename)
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

	// update to tb_staging_ups
	err = service.stagingUpsRepository.UpdateWithTx(camelMap, tx)
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
		Category:          os.Getenv("CATEGORY_STAGING_UPS"),
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

func (service *stagingUpsService) AllSubmittedData() response.UniversalResponse {
	var response response.UniversalResponse
	data := service.stagingUpsRepository.FindAllSubmittedData()
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

func (service *stagingUpsService) GetSubmittedDataBySn(request stagingUpsRequest.FindBySn) response.UniversalResponse {
	var response response.UniversalResponse
	data := service.stagingUpsRepository.FindBySn(request.Sn)
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

func (service *stagingUpsService) GetRejectedData(request stagingUpsRequest.FindRejectedData) response.UniversalResponse {
	var response response.UniversalResponse
	data := service.stagingUpsRepository.FindRejectedData(request.IdUploader)
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
