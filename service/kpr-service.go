package service

import (
	"go-api/dto/request/KPRRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/entity"
	"go-api/repository"
	"time"

	"github.com/mashingan/smapping"
)

type KPRService interface {
	PengajuanKPR(request KPRRequestDTO.PengajuanKPRRequest) responseDTO.Response
}

type kprService struct {
	customerRepository repository.CustomerRepository
	kprRepository      repository.KPRRepository
	salesRepository    repository.SalesRepository
}

func NewKPRService(customerRepo repository.CustomerRepository, kprRepo repository.KPRRepository, salesRepo repository.SalesRepository) KPRService {
	return &kprService{
		customerRepository: customerRepo,
		kprRepository:      kprRepo,
		salesRepository:    salesRepo,
	}
}

func (service *kprService) PengajuanKPR(request KPRRequestDTO.PengajuanKPRRequest) responseDTO.Response {
	var response responseDTO.Response

	customer := entity.TblCustomer{}
	pengajuanKPR := entity.TblPengajuanKprBySales{}
	customer.CreatedAt = time.Now()
	customer.ModifiedAt = time.Now()
	pengajuanKPR.CreatedAt = time.Now()
	pengajuanKPR.ModifiedAt = time.Now()

	salesID := service.salesRepository.GetIDByRefCode(request.ReferralCode)
	request.SalesID = salesID

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

	err = smapping.FillStruct(&pengajuanKPR, smapping.MapFields(&request))
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
	service.kprRepository.Insert(pengajuanKPR)

	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseData = data
	response.ResponseDesc = "Success"
	response.Summary = nil

	return response
}
