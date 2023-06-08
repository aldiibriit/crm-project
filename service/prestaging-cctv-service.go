package service

import (
	"encoding/json"
	logActivityRequest "go-api/dto/request/log-activity"
	prestagingCctvRequest "go-api/dto/request/prestaging-cctv"
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

type PrestagingCCTVService interface {
	PostPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingCctvRequest.PostPrestaging) response.UniversalResponse
	ApprovePrestaging(request prestagingCctvRequest.ApprovePrestaging) response.UniversalResponse
	RejectPrestaging(request prestagingCctvRequest.RejectPrestaging) response.UniversalResponse
	ReuploadPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingCctvRequest.PostPrestaging) response.UniversalResponse
	AllSubmittedData() response.UniversalResponse
	GetSubmittedDataBySn(request prestagingCctvRequest.FindBySn) response.UniversalResponse
	GetRejectedData(request prestagingCctvRequest.FindRejectedData) response.UniversalResponse
}

type prestagingCCTVService struct {
	minioRepository          ExternalRepository.MinioRepository
	logActivityRepository    InternalRepository.LogActivityRepository
	prestagingCCTVRepository InternalRepository.PrestagingCCTVRepository
	baseRepository           InternalRepository.BaseRepository
	stagingCCTVRepository    InternalRepository.StagingCCTVRepository
}

func NewPrestagingCCTVService(minioRepo ExternalRepository.MinioRepository, logActivityRepo InternalRepository.LogActivityRepository, prestagingCCTVRepo InternalRepository.PrestagingCCTVRepository, baseRepo InternalRepository.BaseRepository, stagingCCTVRepo InternalRepository.StagingCCTVRepository) PrestagingCCTVService {
	return &prestagingCCTVService{
		minioRepository:          minioRepo,
		baseRepository:           baseRepo,
		logActivityRepository:    logActivityRepo,
		prestagingCCTVRepository: prestagingCCTVRepo,
		stagingCCTVRepository:    stagingCCTVRepo,
	}
}

func (service *prestagingCCTVService) PostPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingCctvRequest.PostPrestaging) response.UniversalResponse {
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
	tbPrestagingCctv := entity.TbPrestagingCctv{}
	err = smapping.FillStruct(&tbPrestagingCctv, smapping.MapFields(&request))
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("MAPPING_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	// insert to tb_prestaging_DigitalSignage
	err = service.prestagingCCTVRepository.InsertWithTx(tbPrestagingCctv, tx)
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

func (service *prestagingCCTVService) ApprovePrestaging(request prestagingCctvRequest.ApprovePrestaging) response.UniversalResponse {
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

	dataExist := service.stagingCCTVRepository.FindBySn(request.Sn)
	if dataExist.Sn == "" || len(dataExist.Sn) == 0 {
		tbStagingCctv := entity.TbStagingCctv{
			Sn: request.Sn,
		}
		err := service.stagingCCTVRepository.InsertWithTx(tbStagingCctv, tx)
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

func (service *prestagingCCTVService) RejectPrestaging(request prestagingCctvRequest.RejectPrestaging) response.UniversalResponse {
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
	err = service.prestagingCCTVRepository.UpdateWithTx(snakeCaseMap, tx)
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

func (service *prestagingCCTVService) ReuploadPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingCctvRequest.PostPrestaging) response.UniversalResponse {
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
	err = service.prestagingCCTVRepository.UpdateWithTx(camelMap, tx)
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

func (service *prestagingCCTVService) AllSubmittedData() response.UniversalResponse {
	var response response.UniversalResponse
	data := service.prestagingCCTVRepository.FindAllSubmittedData()
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

func (service *prestagingCCTVService) GetSubmittedDataBySn(request prestagingCctvRequest.FindBySn) response.UniversalResponse {
	var response response.UniversalResponse
	data := service.prestagingCCTVRepository.FindBySn(request.Sn)
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

func (service *prestagingCCTVService) GetRejectedData(request prestagingCctvRequest.FindRejectedData) response.UniversalResponse {
	var response response.UniversalResponse
	data := service.prestagingCCTVRepository.FindRejectedData(request.IdUploader)
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
