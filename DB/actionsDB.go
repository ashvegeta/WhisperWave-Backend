package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Table struct {
	Client    *dynamodb.Client
	TableName string
}

func (tableObj *Table) CreateTable(tableName string, tableDef []types.AttributeDefinition, Schema []types.KeySchemaElement) (*types.TableDescription, error) {
	var tableDesc *types.TableDescription

	table, err := tableObj.Client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: tableDef,
		KeySchema:            Schema,
	})

	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", tableObj.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(tableObj.Client)

		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableObj.TableName)}, 5*time.Minute)

		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}

	return tableDesc, err
}

func (tableObj *Table) ListTables() {
    // Build the request with its input parameters
    resp, err := tableObj.Client.ListTables(context.TODO(), &dynamodb.ListTablesInput{
        Limit: aws.Int32(5),
    })
    if err != nil {
        log.Fatalf("failed to list tables, %v", err)
    }

    // list tables
    fmt.Println("Tables:")
    for _, tableName := range resp.TableNames {
        fmt.Println(tableName)
    }
}

func (tableObj *Table) addItemToTable() {

}

func (tableObj *Table) getItemFromTable() {

}

func (tableObj *Table) updateItemInTable() {

}

func (tableObj *Table) deleteItemFromTable() {

}