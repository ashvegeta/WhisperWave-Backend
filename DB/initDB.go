package db

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func getDBClient(cfg aws.Config) *dynamodb.Client {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files

	// Using the Config value, create the DynamoDB client
	return dynamodb.NewFromConfig(cfg)
}

func InitializeTables() {

}