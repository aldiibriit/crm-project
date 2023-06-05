package config

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func SetupMinioConnection() *minio.Client {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "Q76A31bDdwwJJ4nlRvBF"
	secretAccessKey := "ujIJBFY3AO1x4kBgl0xLkDlOOABu2U0Hq917O7WT"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// log.Printf("%#v\n", minioClient)
	return minioClient
}
