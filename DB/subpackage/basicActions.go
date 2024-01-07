package subpackage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// List all dynamoDB tables
func ListTables(client dynamodb.Client, limit int) {
	// Build the request with its input parameters
	resp, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(int32(limit)),
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	// list tables
	fmt.Println("Tables:")
	for c, tableName := range resp.TableNames {
		fmt.Println(c, ". ", tableName)
	}
}

// Create table on dynamoDB
func CreateTable(dbclient *dynamodb.Client, tableInfo *dynamodb.CreateTableInput) (*types.TableDescription, error) {
	var tableDesc *types.TableDescription

	table, err := dbclient.CreateTable(context.TODO(), tableInfo)

	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", tableInfo.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(dbclient)

		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(*tableInfo.TableName)}, 5*time.Minute)

		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}

	return tableDesc, err
}
