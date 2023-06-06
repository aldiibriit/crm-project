package service

import (
	"encoding/json"
	logActivityRequest "go-api/dto/request/log-activity"
	prestagingCRMRequest "go-api/dto/request/prestaging-crm"
	"go-api/dto/response"
	"go-api/entity"
	ExternalRepository "go-api/repository/external-repo"
	InternalRepository "go-api/repository/internal-repo"
	"io"
	"log"
	"mime/multipart"
	"os"
	"reflect"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/mashingan/smapping"
)

type PrestagingCRMService interface {
	PostPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingCRMRequest.PostPrestaging) response.UniversalResponse
	Approve(request prestagingCRMRequest.ApprovePrestaging) response.UniversalResponse
	Reject(request prestagingCRMRequest.RejectPrestaging) response.UniversalResponse
	ReuploadPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingCRMRequest.PostPrestaging) response.UniversalResponse
	AllSubmittedData() response.UniversalResponse
	PostPrestagingV2(request []*multipart.FileHeader) response.UniversalResponse
}

type prestagingCRMService struct {
	minioRepository       ExternalRepository.MinioRepository
	logActivityRepository InternalRepository.LogActivityRepository
	prestagingRepository  InternalRepository.PrestagingCRMRepository
	baseRepository        InternalRepository.BaseRepository
	stagingRepository     InternalRepository.StagingCRMRepository
}

func NewPrestagingCRMService(minioRepo ExternalRepository.MinioRepository, logActivityRepo InternalRepository.LogActivityRepository, prestagingRepo InternalRepository.PrestagingCRMRepository, baseRepo InternalRepository.BaseRepository, stagingRepo InternalRepository.StagingCRMRepository) PrestagingCRMService {
	return &prestagingCRMService{
		minioRepository:       minioRepo,
		baseRepository:        baseRepo,
		logActivityRepository: logActivityRepo,
		prestagingRepository:  prestagingRepo,
		stagingRepository:     stagingRepo,
	}
}

func (service *prestagingCRMService) PostPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingCRMRequest.PostPrestaging) response.UniversalResponse {
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
				err = service.minioRepository.UploadFile(fileReader, v.Size, "crm-project", "prestaging-crm/"+v.Filename)
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
	tbPrestagingCrm := entity.TbPrestagingCrm{}
	err = smapping.FillStruct(&tbPrestagingCrm, smapping.MapFields(&request))
	if err != nil {
		log.Println(err.Error())
		response.HttpCode = 500
		response.ResponseCode = "99"
		response.ResponseMessage = os.Getenv("MAPPING_ERROR_MESSAGE")
		response.Data = nil
		return response
	}

	// insert to tb_prestaging_crm
	err = service.prestagingRepository.InsertWithTx(tbPrestagingCrm, tx)
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
		Category:          os.Getenv("CATEGORY_PRESTAGING_CRM"),
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

func (service *prestagingCRMService) Approve(request prestagingCRMRequest.ApprovePrestaging) response.UniversalResponse {
	var response response.UniversalResponse

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	// insert to tbl log
	tx := service.baseRepository.StartTransaction()
	tblLogActivity := entity.TbLogActivity{
		Category:          os.Getenv("CATEGORY_PRESTAGING_CRM"),
		Sn:                request.Sn,
		StatusDescription: os.Getenv("LOG_APPROVED"),
	}

	dataExist := service.stagingRepository.FindBySn(request.Sn)
	if dataExist.Sn == "" || len(dataExist.Sn) == 0 {
		tbStagingCRM := entity.TbStagingCrm{
			Sn: request.Sn,
		}
		err := service.stagingRepository.InsertWithTx(tbStagingCRM, tx)
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

func (service *prestagingCRMService) Reject(request prestagingCRMRequest.RejectPrestaging) response.UniversalResponse {
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
	err = service.prestagingRepository.UpdateWithTx(snakeCaseMap, tx)
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
		Category:          os.Getenv("CATEGORY_PRESTAGING_CRM"),
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

func (service *prestagingCRMService) ReuploadPrestaging(requestMap map[string]*multipart.FileHeader, request prestagingCRMRequest.PostPrestaging) response.UniversalResponse {
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
				err = service.minioRepository.UploadFile(fileReader, v.Size, "crm-project", "prestaging-crm/"+v.Filename)
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

	// update to tb_prestaging_crm
	err = service.prestagingRepository.UpdateWithTx(camelMap, tx)
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
		Category:          os.Getenv("CATEGORY_PRESTAGING_CRM"),
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

func (service *prestagingCRMService) AllSubmittedData() response.UniversalResponse {
	var response response.UniversalResponse
	data := service.prestagingRepository.FindAllSubmittedData()
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

func (service *prestagingCRMService) PostPrestagingV2(request []*multipart.FileHeader) response.UniversalResponse {
	var response response.UniversalResponse
	go func() {
		// time.Sleep purposed to hold the create file process (expected : the create file function will execute when response sent to front end,so the client doesn't need to wait for file created)

		time.Sleep(5 * time.Second)
		// count the process
		now := time.Now()
		defer timeTrack(now, "createFile")

		for _, v := range request {
			// buat file
			file, err := os.Create("tempFile/" + v.Filename)
			if err != nil {
				log.Println(err.Error())
			}
			defer file.Close()

			srcFile, err := v.Open()
			if err != nil {
				log.Println(err.Error())
			}
			defer srcFile.Close()

			_, err = io.Copy(file, srcFile)
			if err != nil {
				log.Println(err.Error())
			}

			log.Println("file saved")
		}
	}()

	var x prestagingCRMRequest.PostPrestaging

	mappingFilesToStruct(request, x)

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Data = request

	return response
}

func timeTrack(start time.Time, name string) {
	log.Println("the", name, "process take", int(time.Since(start)/time.Second), "second")
}

func mappingFilesToStruct(data []*multipart.FileHeader, result any) {
	log.Println("data type of result", reflect.TypeOf(result))
}
