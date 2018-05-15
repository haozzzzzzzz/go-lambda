package table

import (
	"github.com/haozzzzzzzz/go-lambda/resource/dynamodb"
)

// 简单存储模型
type SimpleStoreModel struct {
	PartitionKey string      `json:"partition_key"`
	SortKey      string      `json:"sort_key"`
	Value        interface{} `json:"value"`
}

type SimpleStoreTable struct {
	dynamodb.DynamoDBTable
}
