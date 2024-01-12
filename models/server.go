package models

type ServerInfo struct {
	SrvName string
	SrvAddr string
	MQ      MessageQueue
}

type UserServerMap struct {
	UserID     string     `dynamodbav:"UserID"`
	ServerInfo ServerInfo `dynamodbav:"ServerInfo"`
}

type MessageQueue struct {
	MQName   string `json:"mq_name"`
	MQURI    string `json:"mq_uri"`
	MQParams []any  `json:"mq_params"`
}
