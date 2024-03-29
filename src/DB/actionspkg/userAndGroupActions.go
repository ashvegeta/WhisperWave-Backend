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

var tableStructUGA *TableStruct

// Initialize the DS for operations
func InitUserAndGroupActions(db_client *dynamodb.Client, tableName string) {
	tableStructUGA = &TableStruct{
		DBClient:  db_client,
		TableName: tableName,
	}
}

// Adds New Group or User to DynamoDB
func AddNewUserOrGroup(newItem any) error {
	// ASSERT
	switch newItem.(type) {
	case models.User:
		break
	case models.Group:
		break
	default:
		return fmt.Errorf("invalid data type (needs to be of type \"User\" or \"Group\")")
	}

	av, err := attributevalue.MarshalMap(newItem)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	_, err = tableStructUGA.DBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableStructUGA.TableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put Record to DynamoDB, %v", err)
	}

	return nil
}

// Gets User Info from DynamoDB (CAN combine logic with below function)
func GetUserInfo(userInfo models.UserOrGroupParams) (models.User, error) {
	var (
		err      error
		response *dynamodb.QueryOutput
		users    []models.User
	)

	// Build Key Expression
	keyEx := expression.Key("ID").Equal(expression.Value(userInfo.PK))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		err = fmt.Errorf("couldn't build expression for query. Here's why: %v", err)

	} else {

		// Query
		response, err = tableStructUGA.DBClient.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(tableStructUGA.TableName),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		})

		if err != nil {
			err = fmt.Errorf("couldn't query for PK : %v. Here's why: %v", userInfo.PK, err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &users)
			if err != nil {
				err = fmt.Errorf("couldn't unmarshal query response. Here's why: %v", err)
			}
		}
	}

	if len(users) == 1 {
		return users[0], err
	}

	return models.User{}, err
}

// Gets Group Info from DynamoDB
func GetGroupInfo(params models.UserOrGroupParams) ([]models.Group, error) {
	var (
		err      error
		response *dynamodb.QueryOutput
		groups   []models.Group
	)

	// Build Key Expression
	keyEx := expression.Key("ID").Equal(expression.Value(params.PK))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		err = fmt.Errorf("couldn't build expression for query. Here's why: %v", err)

	} else {
		// Query
		response, err = tableStructUGA.DBClient.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(tableStructUGA.TableName),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		})

		if err != nil {
			err = fmt.Errorf("couldn't query for PK : %v. Here's why: %v", params.PK, err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &groups)
			if err != nil {
				err = fmt.Errorf("couldn't unmarshal query response. Here's why: %v", err)
			}
		}
	}

	return groups, err
}

// Update a user or group information
func UpdateUserOrGroupInfo(params models.UserOrGroupParams, item any) (interface{}, error) {
	// ASSERT
	switch item.(type) {
	case models.User:
		break
	case models.Group:
		break
	default:
		return nil, fmt.Errorf("invalid data type (needs to be of type \"User\" or \"Group\")")
	}

	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)
	var update expression.UpdateBuilder
	var attributeMap interface{}

	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).IsZero() {
			update = update.Set(expression.Name(t.Field(i).Name), expression.Value(v.Field(i).Interface()))
		}
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		response, err := tableStructUGA.DBClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName:                 aws.String(tableStructUGA.TableName),
			Key:                       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: params.PK}},
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

func UpdateGroupOrUserList(id string, OP string, attrName string, attrValues []string) error {
	// Update the item
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableStructUGA.TableName),
		Key:                       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression:          aws.String(OP + " #attributeName :elementToChange"),
		ExpressionAttributeNames:  map[string]string{"#attributeName": *aws.String(attrName)},
		ExpressionAttributeValues: map[string]types.AttributeValue{":elementToChange": &types.AttributeValueMemberSS{Value: attrValues}},
		ReturnValues:              types.ReturnValueAllNew,
	}

	// perform action
	_, err := tableStructUGA.DBClient.UpdateItem(context.TODO(), updateItemInput)
	if err != nil {
		return err
	}

	return nil
}

// Delete a particular group or user
func DeleteUserOrGroup(params models.UserOrGroupParams) (map[string]types.AttributeValue, error) {

	// Delete user/group
	deletedOutput, err := tableStructUGA.DBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName:    aws.String(tableStructUGA.TableName),
		Key:          map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: params.PK}},
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't delete %v from the table. Here's why: %v", params.PK, err)
	}

	var userInfo models.User
	attributevalue.UnmarshalMap(deletedOutput.Attributes, &userInfo)

	// Cascading deletes
	items, err := LoadChatHistory(models.ChatParams{PK: params.PK})
	if err != nil {
		return nil, fmt.Errorf("couldn't load chat history for cascading deletes for user/group : %s, %s", params.PK, err)
	}

	// 1. Delete chat history of the user/group
	for _, item := range items {
		err := DeleteSingleChat(models.ChatParams{PK: params.PK, SK: item.SK})
		if err != nil {
			return nil, fmt.Errorf("couldn't perform cascading delete (chat history) for user/group : %s, %s", params.PK, err)
		}
	}

	// 2. IF USER - Delete user's group chats (NEED TO USE A GLOBAL INDEX)
	// DeleteUserGroupChat(models.ChatParams{PK: params.PK})

	// 3. IF USER - Delete user's servermap
	err = DeleteServerMap(params.PK)
	if err != nil {
		return nil, err
	}

	// 4. IF USER - Delete user in all group's user list
	for _, gid := range userInfo.GroupList {
		UpdateGroupOrUserList(gid, "DELETE", "UserList", []string{userInfo.UserId})
	}

	// 5. IF USER - Delete user id from user's friend's friend list
	for _, uid := range userInfo.FriendsList {
		UpdateGroupOrUserList(uid, "DELETE", "FriendsList", []string{userInfo.UserId})
	}

	return deletedOutput.Attributes, err
}
