package revenueModel

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RevenueModel struct {
	db     *gorm.DB
	client *minio.Client
}

func New(dsn, endpoint, accessKey, secretKey string, useSSL bool) (*RevenueModel, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &RevenueModel{
		db:     db,
		client: client,
	}, nil
}

func NewRevenueModel() (*RevenueModel, error) {
	return &RevenueModel{}, nil
}
