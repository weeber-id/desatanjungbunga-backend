package services

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
)

// MinioClient session variable
var MinioClient *minio.Client

// InitializationMinio service for object storage
func InitializationMinio() {
	config := variables.MinioConfig

	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Cannot connect to Minio Server: %v", err)
	}

	MinioClient = client
}
