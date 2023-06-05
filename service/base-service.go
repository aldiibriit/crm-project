package service

import (
	InternalRepository "go-api/repository/internal-repo"

	"gorm.io/gorm"
)

type BaseService interface {
	StartTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB)
	RollbackTransaction(tx *gorm.DB)
}

type baseService struct {
	baseRepository InternalRepository.BaseRepository
}

func NewBaseService(baseRepo InternalRepository.BaseRepository) BaseService {
	return &baseService{
		baseRepository: baseRepo,
	}
}

func (service *baseService) StartTransaction() *gorm.DB {
	return service.baseRepository.StartTransaction()
}

func (service *baseService) CommitTransaction(tx *gorm.DB) {
	service.baseRepository.CommitTransaction(tx)
}

func (service *baseService) RollbackTransaction(tx *gorm.DB) {
	service.baseRepository.RollbackTransaction(tx)
}
