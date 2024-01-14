package subpackage

import (
	"WhisperWave-BackEnd/src/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

type TableStruct struct {
	DBClient  *dynamodb.Client
	TableName string
}

func LoadDefaultConfig() *dynamodb.Client {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println(err)
		return nil
	}

	return GetDBClient(config)
}

// Using the Config value, create the DynamoDB client
func GetDBClient(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

// Initialize all the tables by loading schema
func InitializeTables(db_client *dynamodb.Client) {

	// load dynamoDB config
	var dbConfig models.DBConfig

	file, err := os.Open("./config/db.json")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dbConfig)
	if err != nil {
		log.Panicln(err)
		return
	}

	// iterate over tables and create dynamoDB tables
	var tables []*types.TableDescription

	for _, tableInfo := range dbConfig.Tables {
		// Create table description
		tableDesc := &dynamodb.CreateTableInput{
			TableName:             &tableInfo.TableName,
			AttributeDefinitions:  tableInfo.Attributes,
			KeySchema:             tableInfo.KeySchema,
			ProvisionedThroughput: &tableInfo.ProvisionedThroughput,
		}

		if len(tableInfo.GSI) > 0 {
			tableDesc.GlobalSecondaryIndexes = tableInfo.GSI
		}

		// create table
		var OPerror *smithy.OperationError
		table, err := CreateTable(db_client, tableDesc)

		if errors.As(err, &OPerror) && OPerror.OperationName == "CreateTable" {
			log.Printf("table \"%s\" already exists .... skipping table creation\n", tableInfo.TableName)
			continue
		} else if err != nil {
			log.Println(err)
			return
		} else {
			log.Printf("table \"%s\" created successfully\n", tableInfo.TableName)
			tables = append(tables, table)
		}

	}

	if len(tables) > 0 {
		fmt.Print(tables)
	}

	// Init Table Structs
	InitChatHistory(db_client, "ChatHistory")
	InitUserAndGroupActions(db_client, "UserAndGroupInfo")
	InitUserServerMap(db_client, "UserServerMap")
}
