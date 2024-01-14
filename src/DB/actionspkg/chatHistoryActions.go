package subpackage

import (
	"WhisperWave-BackEnd/src/models"
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var tableStructCH *TableStruct

// Initialize the DS for operations
func InitChatHistory(db_client *dynamodb.Client, tableName string) {
	tableStructCH = &TableStruct{
		DBClient:  db_client,
		TableName: tableName,
	}
}

// Add new chat text b/w 2 users or a group
func AddNewChat(chat models.ChatHistory) error {
	av, err := attributevalue.MarshalMap(chat)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}

	_, err = tableStructCH.DBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableStructCH.TableName),
		Item:      av,
	})

	return err
}

// Load the chat history b/w two users or a group
func LoadChatHistory(chat models.ChatParams) ([]models.ChatHistory, error) {
	var (
		err         error
		response    *dynamodb.QueryOutput
		chatHistory []models.ChatHistory
	)

	// Build Key Expression
	keyEx := expression.Key("ID").Equal(expression.Value(chat.PK))
	if chat.SK != "" {
		keyEx = keyEx.And(expression.Key("UserID-TimeStamp").BeginsWith(chat.SK))
	}

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)

	} else {
		// Query
		response, err = tableStructCH.DBClient.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(tableStructCH.TableName),
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

// Update chat text
func UpdateChat(chat models.ChatParams, item models.ChatHistory) (interface{}, error) {
	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)
	eb := expression.NewBuilder()
	var update expression.UpdateBuilder
	var attributeMap interface{}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name == "Content" {
			update = update.Set(expression.Name(t.Field(i).Name), expression.Value(v.Field(i).Interface()))
		}
	}

	expr, err := eb.WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		response, err := tableStructCH.DBClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName: aws.String(tableStructCH.TableName),
			Key: map[string]types.AttributeValue{
				"ID":               &types.AttributeValueMemberS{Value: chat.PK},
				"UserID-TimeStamp": &types.AttributeValueMemberS{Value: chat.SK},
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			return nil, fmt.Errorf("couldn't update \"User\" OR \"Group\" %v. Here's why: %v", item, err)
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				return nil, fmt.Errorf("couldn't unmarshall update response. Here's why: %v", err)
			}
		}
	}

	return attributeMap, err
}

// Delete chat history
func DeleteSingleChat(chat models.ChatParams) error {
	// set delete input
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableStructCH.TableName),
		Key: map[string]types.AttributeValue{
			"ID":               &types.AttributeValueMemberS{Value: chat.PK},
			"UserID-TimeStamp": &types.AttributeValueMemberS{Value: chat.SK},
		},
	}

	// perform delete operation
	_, err := tableStructCH.DBClient.DeleteItem(context.TODO(), deleteInput)
	if err != nil {
		return fmt.Errorf("couldn't delete the chat involving user : %v from the table. Here's why: %v", chat.PK, err)
	}
	return nil
}

// Delete user's group chat (INCOMPLETE)
func DeleteUserGroupChat(chat models.ChatParams) error {
	user, err := GetUserInfo(models.UserOrGroupParams{PK: chat.PK})
	if err != nil {
		return fmt.Errorf("error in fetching user info : %s", chat.PK)
	}

	for _, gid := range user.GroupList {
		err := DeleteSingleChat(models.ChatParams{
			PK: gid,
			SK: user.UserId,
		})
		if err != nil {
			return fmt.Errorf("error in deleting user : %s, group : %s chat", user.UserId, gid)
		}
	}

	return nil
}
