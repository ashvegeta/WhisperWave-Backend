package models

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TableStruct struct {
	DBClient  *dynamodb.Client
	TableName string
}

type TableInfo struct {
	TableName             string
	Attributes            []types.AttributeDefinition
	KeySchema             []types.KeySchemaElement
	ProvisionedThroughput types.ProvisionedThroughput
	GSI                   []types.GlobalSecondaryIndex
}

type DBConfig struct {
	Tables []TableInfo `json:"Tables"`
}

type ChatParams struct {
	PK string `dynamodbav:"ID"`
	SK string `dynamodbav:"UserID-TimeStamp"`
}

type ChatHistory struct {
	PK      string `dynamodbav:"ID"`
	SK      string `dynamodbav:"UserID-TimeStamp"`
	MID     string `dynamodbav:"MessageID"`
	MType   string `dynamodbav:"MessageType"`
	Content string `dynamodbav:"Content"`
}

type UserOrGroupParams struct {
	PK string `dynamodbav:"ID"`
}

type ServerInfo struct {
	SrvName string
	SrvAddr string
	MQ      MessageQueue
}

type UserServerMap struct {
	UserID     string     `dynamodbav:"UserID"`
	ServerInfo ServerInfo `dynamodbav:"ServerInfo"`
}
