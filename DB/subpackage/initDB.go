package subpackage

import (
	"WhisperWave-BackEnd/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Using the Config value, create the DynamoDB client
func GetDBClient(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

// Initialize all the tables by loading schema
func InitializeTables() {
	// load AWS credentials config
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println(err)
		return
	}

	// load dynamoDB config
	var dbConfig models.DBConfig

	file, err := os.Open("../config/db.json")
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
	dbclient := GetDBClient(config)

	for _, tableInfo := range dbConfig.Tables {
		// extract table metadata
		tableDesc := &dynamodb.CreateTableInput{
			TableName: &tableInfo.TableName,

			AttributeDefinitions: []types.AttributeDefinition{{
				AttributeName: aws.String(tableInfo.AttributeDef.AttributeName),
				AttributeType: types.ScalarAttributeType(tableInfo.AttributeDef.AttributeType),
			}},

			KeySchema: []types.KeySchemaElement{{
				AttributeName: aws.String(tableInfo.KeySchema.AttributeName),
				KeyType:       types.KeyType(tableInfo.KeySchema.KeyType),
			}},

			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  &tableInfo.ProvThroughput.RCU,
				WriteCapacityUnits: &tableInfo.ProvThroughput.WCU,
			},
		}

		// create table
		table, err := CreateTable(dbclient, tableDesc)
		if err != nil {
			log.Println(err)
			return
		}

		tables = append(tables, table)
	}

	fmt.Print(tables)
}
