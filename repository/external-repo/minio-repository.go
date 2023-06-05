package ExternalRepository

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
)

type MinioRepository interface {
	UploadFile(fileReader io.Reader, fileSize int64, bucketName string, objectName string) error
}

type minioConnection struct {
	connection *minio.Client
}

func NewMinioRepository(conn *minio.Client) MinioRepository {
	return &minioConnection{
		connection: conn,
	}
}

func (mc *minioConnection) UploadFile(fileReader io.Reader, fileSize int64, bucketName string, objectName string) error {

	uploadInfo, err := mc.connection.PutObject(context.Background(), bucketName, objectName, fileReader, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("Successfully uploaded", uploadInfo)
	return nil
}
