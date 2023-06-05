package service

import (
	CCTVRequestDTO "go-api/dto/request/cctv"
	"go-api/dto/response"
	internalRepository "go-api/repository/internal-repo"
)

type CCTVService interface {
	GetAll() response.UniversalResponse
	FindBySN(request CCTVRequestDTO.FindBySNRequest) response.UniversalResponse
}

type cctvService struct {
	cctvRepository internalRepository.CCTVRepository
}

func NewCCTVService(cctvRepo internalRepository.CCTVRepository) CCTVService {
	return &cctvService{
		cctvRepository: cctvRepo,
	}
}

func (service *cctvService) GetAll() response.UniversalResponse {
	var response response.UniversalResponse

	cctvs := service.cctvRepository.All()
	if len(cctvs) == 0 {
		response.HttpCode = 404
		response.ResponseCode = "99"
		response.ResponseMessage = "data not found!"
		response.Data = nil
	} else {
		response.HttpCode = 200
		response.ResponseCode = "00"
		response.ResponseMessage = "success"
		response.Data = cctvs
	}
	return response
}

func (service *cctvService) FindBySN(request CCTVRequestDTO.FindBySNRequest) response.UniversalResponse {
	var response response.UniversalResponse
	cctv := service.cctvRepository.FindBySN(request.SnCctv)
	if cctv.SnCctv == "" || len(cctv.SnCctv) == 0 {
		response.HttpCode = 404
		response.ResponseCode = "99"
		response.ResponseMessage = "data not found!"
		response.Data = nil
	} else {
		response.HttpCode = 200
		response.ResponseCode = "00"
		response.ResponseMessage = "success"
		response.Data = cctv
	}
	return response
}
