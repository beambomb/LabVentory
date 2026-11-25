package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"labventroy/backend/internal/config"
	"labventroy/backend/internal/handler"
	"labventroy/backend/internal/repository"
	"labventroy/backend/internal/route"
	"labventroy/backend/internal/service"
)

func main() {
	cfg := config.Load()

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	ddbClient := dynamodb.NewFromConfig(awsCfg)
	repo := repository.NewDynamoDBRepository(ddbClient, cfg.TableName)
	svc := service.NewInventoryService(repo)
	itemHandler := handler.NewItemHandler(svc)
	router := route.NewRouter(itemHandler)

	lambda.Start(router)
}
