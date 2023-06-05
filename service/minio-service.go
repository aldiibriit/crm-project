package service

import (
	externalRepository "go-api/repository/external-repo"
	"io"
	"log"
	"mime/multipart"
)

type MinioService interface {
	UploadFile(file *multipart.FileHeader, bucketName string, objectName string) error
}

type minioService struct {
	minioRepository externalRepository.MinioRepository
}

func NewMinioService(minioRepo externalRepository.MinioRepository) MinioService {
	return &minioService{
		minioRepository: minioRepo,
	}
}

func (service *minioService) UploadFile(file *multipart.FileHeader, bucketName string, objectName string) error {
	openedFile, err := file.Open()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer openedFile.Close()

	fileReader := io.Reader(openedFile)

	service.minioRepository.UploadFile(fileReader, file.Size, bucketName, objectName)

	return nil
}
