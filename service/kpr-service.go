package service

import (
	"go-api/dto/request/kprRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/entity"
	"go-api/repository"
	"time"

	"github.com/mashingan/smapping"
)

type KPRService interface {
	PengajuanKPR(request kprRequestDTO.PengajuanKPRRequest) responseDTO.Response
}

type kprService struct {
	customerRepository repository.CustomerRepository
}

func NewKPRService(customerRepo repository.CustomerRepository) KPRService {
	return &kprService{
		customerRepository: customerRepo,
	}
}

func (service *kprService) PengajuanKPR(request kprRequestDTO.PengajuanKPRRequest) responseDTO.Response {
	var response responseDTO.Response

	customer := entity.TblCustomer{}
	customer.CreatedAt = time.Now()
	err := smapping.FillStruct(&customer, smapping.MapFields(&request))
	if err != nil {
		response.HttpCode = 200
		response.MetadataResponse = nil
		response.ResponseCode = "00"
		response.ResponseData = nil
		response.ResponseDesc = "failed map " + err.Error()
		response.Summary = nil

		return response
	}

	data := service.customerRepository.Insert(customer)
	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = data
	response.ResponseDesc = "Success"
	response.Summary = nil

	return response
}
