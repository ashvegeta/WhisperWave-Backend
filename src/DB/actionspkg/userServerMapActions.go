package subpackage

import (
	"WhisperWave-BackEnd/src/models"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var tableStructUSM *TableStruct

// Initialize the DS for operations
func InitServerMap(db_client *dynamodb.Client, tableName string) {
	tableStructUSM = &TableStruct{
		DBClient:  db_client,
		TableName: tableName,
	}
}

// Get the websocket connection map
func GetServerMap(userID string) ([]models.UserServerMap, error) {
	var (
		items    []models.UserServerMap
		err      error
		response *dynamodb.QueryOutput
	)

	// Build Key Expression
	keyExp := expression.Key("UserID").Equal(expression.Value(userID))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExp).Build()

	if err != nil {
		err = fmt.Errorf("couldn't build expression for query. Here's why: %v", err)

	} else {
		// Query
		response, err = tableStructUSM.DBClient.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(tableStructUSM.TableName),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		})

		if err != nil {
			err = fmt.Errorf("couldn't query for PK : %v. Here's why: %v", userID, err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &items)
			if err != nil {
				err = fmt.Errorf("couldn't unmarshal query response. Here's why: %v", err)
			}
		}
	}

	return items, err
}

// Update or Add the websocket connection map
func PutServerMap(usrSrvMap models.UserServerMap) error {
	av, err := attributevalue.MarshalMap(usrSrvMap)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	_, err = tableStructUSM.DBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableStructUSM.TableName),
		Item:      av,
	})

	return err
}

// Delete the websocket connection map
func DeleteServerMap(userID string) error {
	_, err := tableStructUSM.DBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableStructUSM.TableName),
		Key:       map[string]types.AttributeValue{"UserID": &types.AttributeValueMemberS{Value: userID}},
	})
	if err != nil {
		log.Printf("Couldn't delete %v from the table. Here's why: %v\n", userID, err)
	}
	return err
}
