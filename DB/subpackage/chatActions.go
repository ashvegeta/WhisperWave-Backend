package subpackage

import (
	"WhisperWave-BackEnd/models"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func AddNewChat(db_client *dynamodb.Client, tableName string, chat models.ChatHistory) error {
	av, err := attributevalue.MarshalMap(chat)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}

	_, err = db_client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	return err
}

func LoadChatHistory(db_client *dynamodb.Client, tableName string, chat models.ChatParams) ([]models.ChatHistory, error) {
	var (
		err         error
		response    *dynamodb.QueryOutput
		chatHistory []models.ChatHistory
	)

	// Build Key Expression
	keyEx := expression.Key("ID").Equal(expression.Value(chat.PK)).
		And(expression.Key("UserID-TimeStamp").BeginsWith(chat.SK))

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)

	} else {
		// Query
		response, err = db_client.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		})

		if err != nil {
			log.Printf("Couldn't query for PK : %v. Here's why: %v\n", chat.PK, err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &chatHistory)
			if err != nil {
				log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
			}
		}
	}

	return chatHistory, err
}

func DeleteChat() error {
	return fmt.Errorf("")
}
