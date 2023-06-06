package service

import (
	logActivityRequest "go-api/dto/request/log-activity"
	LogActivityResponseDTO "go-api/dto/response/log-activity"
	"go-api/entity"
	internalRepo "go-api/repository/internal-repo"
	"log"
	"strings"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type LogActivityService interface {
	Insert(logActivityRequest.InsertRequest) error
	InsertWithTx(request logActivityRequest.InsertRequest, tx *gorm.DB) error
	GetTimeLine(request logActivityRequest.GetTimeLine) LogActivityResponseDTO.Response
}

type logActivityService struct {
	logActivityRepository internalRepo.LogActivityRepository
}

func NewLogActivityService(logActivityRepo internalRepo.LogActivityRepository) LogActivityService {
	return &logActivityService{
		logActivityRepository: logActivityRepo,
	}
}

func (service *logActivityService) Insert(request logActivityRequest.InsertRequest) error {
	tbActivityLog := entity.TbLogActivity{}
	err := smapping.FillStruct(&tbActivityLog, smapping.MapFields(&request))
	if err != nil {
		log.Println(err.Error())
	}

	err = service.logActivityRepository.Insert(tbActivityLog)
	return err
}

func (service *logActivityService) InsertWithTx(request logActivityRequest.InsertRequest, tx *gorm.DB) error {
	tbActivityLog := entity.TbLogActivity{}
	err := smapping.FillStruct(&tbActivityLog, smapping.MapFields(&request))
	if err != nil {
		log.Println(err.Error())
	}

	err = service.logActivityRepository.InsertWithTx(tbActivityLog, tx)
	return err
}

func (service *logActivityService) GetTimeLine(request logActivityRequest.GetTimeLine) LogActivityResponseDTO.Response {
	var response LogActivityResponseDTO.Response
	var timeline LogActivityResponseDTO.Timeline
	sliceOfTimeline := service.logActivityRepository.FindTimeLineBySn(request.Sn)
	log.Println(sliceOfTimeline)

	for _, v := range sliceOfTimeline {

		if strings.HasPrefix(v.Category, "PRESTAGING") {
			timeline.Prestaging = v
		}

		if strings.HasPrefix(v.Category, "STAGING") && !strings.HasPrefix(v.Category, "STAGING_LIVE") {
			timeline.Staging = v
		}

		if strings.HasPrefix(v.Category, "STAGING_LIVE") {
			timeline.StagingLive = v
		}
	}

	response.HttpCode = 200
	response.ResponseCode = "00"
	response.ResponseMessage = "Success"
	response.Timeline = timeline

	return response
}
