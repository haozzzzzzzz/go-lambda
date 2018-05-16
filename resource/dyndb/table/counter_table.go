package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// 计数器
type CounterModel struct {
	PartitionKey string `form:"partition_key"`
	SortKey      string `form:"sort_key"`
	Count        uint32 `json:"count"`
}

type CounterTable struct {
	DynamoDBTable
}

func (m *CounterTable) Incr(partitionKey string, sortKey string, incrNum uint32) (newNum uint32, err error) {
	return m.IncrCounter(map[string]*dynamodb.AttributeValue{
		"partition_key": {
			S: aws.String(partitionKey),
		},
		"sort_key": {
			S: aws.String(sortKey),
		},
	}, "count", incrNum)
}
