package service

import (
	"go-api/dto/response"
	internalRepository "go-api/repository/internal-repo"
)

type CRMService interface {
	GetAll() response.UniversalResponse
}

type crmService struct {
	crmRepository internalRepository.CRMRepository
}

func NewCRMService(crmRepo internalRepository.CRMRepository) CRMService {
	return &crmService{
		crmRepository: crmRepo,
	}
}

func (service *crmService) GetAll() response.UniversalResponse {
	var response response.UniversalResponse
	data := service.crmRepository.GetAll()
	response.ResponseCode = "00"
	response.ResponseMessage = "success"
	response.Data = data
	return response
}
