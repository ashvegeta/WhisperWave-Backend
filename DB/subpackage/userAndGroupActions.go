package subpackage

import (
	"WhisperWave-BackEnd/models"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Adds New User to the DynamoDB (CAN combine logic with below function)
func AddNewUser(db_client *dynamodb.Client, tableName string, userInfo models.User) error {
	av, err := attributevalue.MarshalMap(userInfo)
	if err != nil {
		return fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	_, err = db_client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put Record to DynamoDB, %v", err)
	}

	return err
}

// Adds New Group Info to the DynamoDB
func AddNewGroup(db_client *dynamodb.Client, tableName string, groupInfo models.Group) error {
	var returnErr error

	av, err := attributevalue.MarshalMap(groupInfo)
	if err != nil {
		returnErr = fmt.Errorf("failed to DynamoDB marshal Record, %v", err)
	}

	_, err = db_client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		returnErr = fmt.Errorf("failed to put Record to DynamoDB, %v", err)
	}

	return returnErr
}

// Gets User Info from DynamoDB (CAN combine logic with below function)
func GetUserInfo(db_client *dynamodb.Client, tableName string, userInfo models.UserOrGroupParams) ([]models.User, error) {
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
		response, err = db_client.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
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

	return users, err
}

// Gets Group Info from DynamoDB
func GetGroupInfo(db_client *dynamodb.Client, tableName string, userInfo models.UserOrGroupParams) ([]models.Group, error) {
	var (
		err      error
		response *dynamodb.QueryOutput
		groups   []models.Group
	)

	// Build Key Expression
	keyEx := expression.Key("ID").Equal(expression.Value(userInfo.PK))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		err = fmt.Errorf("couldn't build expression for query. Here's why: %v", err)

	} else {
		// Query
		response, err = db_client.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		})

		if err != nil {
			err = fmt.Errorf("couldn't query for PK : %v. Here's why: %v", userInfo.PK, err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &groups)
			if err != nil {
				err = fmt.Errorf("couldn't unmarshal query response. Here's why: %v", err)
			}
		}
	}

	return groups, err
}

func UpdateUserInfo() error {
	return fmt.Errorf("")
}

func UpdateGroupInfo() error {
	return fmt.Errorf("")
}

func DeleteUser() error {
	return fmt.Errorf("")
}

func DeleteGroup() error {
	return fmt.Errorf("")
}
