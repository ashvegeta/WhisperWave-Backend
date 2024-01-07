package subpackage

import (
	"WhisperWave-BackEnd/models"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Using the Config value, create the DynamoDB client
func GetDBClient(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

// Initialize all the tables by loading schema
func InitializeTables(db_client *dynamodb.Client) {

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

	for _, tableInfo := range dbConfig.Tables {
		// Create table description
		tableDesc := &dynamodb.CreateTableInput{
			TableName: &tableInfo.TableName,

			AttributeDefinitions: tableInfo.Attributes,

			KeySchema: tableInfo.KeySchema,

			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  tableInfo.ProvisionedThroughput.ReadCapacityUnits,
				WriteCapacityUnits: tableInfo.ProvisionedThroughput.WriteCapacityUnits,
			},
		}

		if len(tableInfo.GSI) > 0 {
			tableDesc.GlobalSecondaryIndexes = tableInfo.GSI
		}

		// create table
		table, err := CreateTable(db_client, tableDesc)
		if err != nil {
			log.Println(err)
			return
		}

		tables = append(tables, table)
	}

	fmt.Print(tables)
}
