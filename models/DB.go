package models

type TableInfo struct {
	TableName      string         `json:"tableName"`
	AttributeDef   AttributeType  `json:"attrDef"`
	KeySchema      KeySchema      `json:"keySchema"`
	ProvThroughput ProvThroughput `json:"provThroughput"`
}

type AttributeType struct {
	AttributeName string `json:"attrName"`
	AttributeType string `json:"attrType"`
}

type KeySchema struct {
	AttributeName string `json:"attrName"`
	KeyType       string `json:"keyType"`
}

type ProvThroughput struct {
	RCU int64 `json:"RCU"`
	WCU int64 `json:"WCU"`
}

type DBConfig struct {
	Tables []TableInfo `json:"tables"`
}
