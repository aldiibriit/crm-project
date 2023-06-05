package service

import (
	logActivityRequest "go-api/dto/request/log-activity"
	"go-api/entity"
	internalRepo "go-api/repository/internal-repo"
	"log"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type LogActivityService interface {
	Insert(logActivityRequest.InsertRequest) error
	InsertWithTx(request logActivityRequest.InsertRequest, tx *gorm.DB) error
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
