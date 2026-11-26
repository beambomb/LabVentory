package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"labventroy/backend/internal/config"
	"labventroy/backend/internal/handler"
	"labventroy/backend/internal/repository"
	"labventroy/backend/internal/route"
	"labventroy/backend/internal/service"
	"labventroy/backend/internal/storage"
)

func main() {
	cfg := config.Load()

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	ddbClient := dynamodb.NewFromConfig(awsCfg)
	repo := repository.NewDynamoDBRepository(ddbClient, cfg.TableName)

	var media storage.MediaStorage = storage.NoopMediaStorage{}
	if cfg.ImagesBucket != "" {
		s3Client := s3.NewFromConfig(awsCfg)
		media = storage.NewS3MediaStorage(s3Client, cfg.ImagesBucket)
	}

	svc := service.NewInventoryServiceWithStorage(repo, media)
	itemHandler := handler.NewItemHandler(svc)
	router := route.NewRouter(itemHandler)

	lambda.Start(router)
}
