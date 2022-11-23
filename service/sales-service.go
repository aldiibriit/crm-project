package service

import (
	"go-api/dto/request/salesRequestDTO"
	responseDTO "go-api/dto/response"
	"go-api/dto/response/salesResponseDTO"
	"go-api/entity"
	"go-api/helper"
	"go-api/repository"
	"strconv"
)

type SalesService interface {
	MISDeveloper(request salesRequestDTO.AllRequest) responseDTO.Response
	MISSuperAdmin() responseDTO.Response
}

type salesService struct {
	salesRepository repository.SalesRepository
}

func NewSalesService(salesRepo repository.SalesRepository) SalesService {
	return &salesService{
		salesRepository: salesRepo,
	}
}

func (service *salesService) MISDeveloper(request salesRequestDTO.AllRequest) responseDTO.Response {
	var response responseDTO.Response

	data := service.salesRepository.FindByEmailDeveloper(request.EmailDeveloper)

	encryptedData := serializeMisDeveloper(data)
	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = encryptedData
	response.Summary = nil

	return response
}

func (service *salesService) MISSuperAdmin() responseDTO.Response {
	var response responseDTO.Response

	data := service.salesRepository.MISSuperAdmin()

	encryptedData := serializeMisSuperAdmin(data)
	response.HttpCode = 200
	response.MetadataResponse = nil
	response.ResponseCode = "00"
	response.ResponseDesc = "Success"
	response.ResponseData = encryptedData
	response.Summary = nil

	return response
}

func serializeMisDeveloper(request interface{}) []entity.TblSales {
	data := request.([]entity.TblSales)
	result := make([]entity.TblSales, len(data))
	for i, v := range data {
		encryptedIdRes, _ := helper.RsaEncryptBEToFE([]byte(strconv.Itoa(v.ID)))
		encryptedEmaiDeveloper, _ := helper.RsaEncryptBEToFE([]byte(v.EmailDeveloper))
		encryptedEmaiSales, _ := helper.RsaEncryptBEToFE([]byte(v.EmailSales))
		encryptedRefferalCode, _ := helper.RsaEncryptBEToFE([]byte(v.RefferalCode))
		encryptedRegisteredBy, _ := helper.RsaEncryptBEToFE([]byte(v.RegisteredBy))
		encryptedSalesName, _ := helper.RsaEncryptBEToFE([]byte(v.SalesName))
		encryptedCreatedAt, _ := helper.RsaEncryptBEToFE([]byte(v.CreatedAt.String()))
		encryptedModifiedAt, _ := helper.RsaEncryptBEToFE([]byte(v.ModifiedAt.String()))

		result[i].IDResponse = encryptedIdRes
		result[i].EmailSales = encryptedEmaiSales
		result[i].EmailDeveloper = encryptedEmaiDeveloper
		result[i].RefferalCode = encryptedRefferalCode
		result[i].RegisteredBy = encryptedRegisteredBy
		result[i].CreatedAtRes = encryptedCreatedAt
		result[i].ModifiedAtRes = encryptedModifiedAt
		result[i].SalesName = encryptedSalesName
	}

	return result
}

func serializeMisSuperAdmin(request interface{}) []salesResponseDTO.MISSuperAdmin {
	data := request.([]salesResponseDTO.MISSuperAdmin)
	result := make([]salesResponseDTO.MISSuperAdmin, len(data))
	for i, v := range data {
		encryptedSalesName, _ := helper.RsaEncryptBEToFE([]byte(v.SalesName))
		encryptedMetadata, _ := helper.RsaEncryptBEToFE([]byte(v.Metadata))
		encryptedStatus, _ := helper.RsaEncryptBEToFE([]byte(v.Status))
		encryptedJenisProperti, _ := helper.RsaEncryptBEToFE([]byte(v.JenisProperti))
		encryptedTipeProperti, _ := helper.RsaEncryptBEToFE([]byte(v.TipeProperti))

		result[i].SalesName = encryptedSalesName
		result[i].Metadata = encryptedMetadata
		result[i].Status = encryptedStatus
		result[i].JenisProperti = encryptedJenisProperti
		result[i].TipeProperti = encryptedTipeProperti
	}

	return result
}
