package service

import (
	"encoding/json"
	logActivityRequest "go-api/dto/request/log-activity"
	stagingCrmRequest "go-api/dto/request/staging-crm"
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

type StagingCRMService interface {
	PostStaging(requestMap map[string]*multipart.FileHeader, request stagingCrmRequest.PostStaging) response.UniversalResponse
	ApproveStaging(request stagingCrmRequest.ApproveStaging) response.UniversalResponse
	RejectStaging(request stagingCrmRequest.RejectStaging) response.UniversalResponse
	ReuploadStaging(requestMap map[string]*multipart.FileHeader, request stagingCrmRequest.PostStaging) response.UniversalResponse
	AllSubmittedData() response.UniversalResponse
	GetSubmittedDataBySn(request stagingCrmRequest.FindBySn) response.UniversalResponse
	GetRejectedData(request stagingCrmRequest.FindRejectedData) response.UniversalResponse
}

type stagingCRMService struct {
	minioRepository          ExternalRepository.MinioRepository
	logActivityRepository    InternalRepository.LogActivityRepository
	baseRepository           InternalRepository.BaseRepository
	stagingCrmRepository     InternalRepository.StagingCRMRepository
	stagingLiveCrmRepository InternalRepository.StagingLiveCRMRepository
}

func NewStagingCRMService(minioRepo ExternalRepository.MinioRepository, logActivityRepo InternalRepository.LogActivityRepository, baseRepo InternalRepository.BaseRepository, stagingCrmRepo InternalRepository.StagingCRMRepository, stagingLiveCrmRepo InternalRepository.StagingLiveCRMRepository) StagingCRMService {
	return &stagingCRMService{
		minioRepository:          minioRepo,
		baseRepository:           baseRepo,
		logActivityRepository:    logActivityRepo,
		stagingCrmRepository:     stagingCrmRepo,
		stagingLiveCrmRepository: stagingLiveCrmRepo,
	}
}

func (service *stagingCRMService) PostStaging(requestMap map[string]*multipart.FileHeader, request stagingCrmRequest.PostStaging) response.UniversalResponse {
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
				err = service.minioRepository.UploadFile(fileReader, v.Size, "crm-project", "staging-crm/"+v.Filename)
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

	// log.Println(updatedData)

	// insert to tb_staging_crm
	dataExist := service.stagingCrmRepository.FindBySn(request.Sn)
	if dataExist.Sn == "" || len(dataExist.Sn) == 0 {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("DATA_NOT_FOUND_MESSAGE")
		response.Data = nil
		return response
	}

	err = service.stagingCrmRepository.UpdateWithTx(snakeCaseMap, tx)
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
		Category:          os.Getenv("CATEGORY_STAGING_CRM"),
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

func (service *stagingCRMService) ApproveStaging(request stagingCrmRequest.ApproveStaging) response.UniversalResponse {
	var response response.UniversalResponse

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	// insert to tbl log
	tx := service.baseRepository.StartTransaction()
	tblLogActivity := entity.TbLogActivity{
		Category:          os.Getenv("CATEGORY_STAGING_CRM"),
		Sn:                request.Sn,
		StatusDescription: os.Getenv("LOG_APPROVED"),
	}

	dataExist := service.stagingLiveCrmRepository.FindBySn(request.Sn)
	if dataExist.Sn == "" || len(dataExist.Sn) == 0 {
		tbStagingLiveCRM := entity.TbStagingLiveCrm{
			Sn: request.Sn,
		}
		err := service.stagingLiveCrmRepository.InsertWithTx(tbStagingLiveCRM, tx)
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

func (service *stagingCRMService) RejectStaging(request stagingCrmRequest.RejectStaging) response.UniversalResponse {
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
	err = service.stagingCrmRepository.UpdateWithTx(snakeCaseMap, tx)
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
		Category:          os.Getenv("CATEGORY_STAGING_CRM"),
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

func (service *stagingCRMService) ReuploadStaging(requestMap map[string]*multipart.FileHeader, request stagingCrmRequest.PostStaging) response.UniversalResponse {
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
				err = service.minioRepository.UploadFile(fileReader, v.Size, "crm-project", "staging-crm/"+v.Filename)
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

	// update to tb_staging_crm
	err = service.stagingCrmRepository.UpdateWithTx(camelMap, tx)
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
		Category:          os.Getenv("CATEGORY_STAGING_CRM"),
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

func (service *stagingCRMService) AllSubmittedData() response.UniversalResponse {
	var response response.UniversalResponse
	data := service.stagingCrmRepository.FindAllSubmittedData()
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

func (service *stagingCRMService) GetSubmittedDataBySn(request stagingCrmRequest.FindBySn) response.UniversalResponse {
	var response response.UniversalResponse
	data := service.stagingCrmRepository.FindBySn(request.Sn)
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

func (service *stagingCRMService) GetRejectedData(request stagingCrmRequest.FindRejectedData) response.UniversalResponse {
	var response response.UniversalResponse
	data := service.stagingCrmRepository.FindRejectedData(request.IdUploader)
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
