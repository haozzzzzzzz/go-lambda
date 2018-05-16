package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// 计数器
type CounterModel struct {
	Name  string `json:"name"`
	Count uint32 `json:"count"`
}

type CounterTable struct {
	DynamoDBTable
}

func (m *CounterTable) Incr(name string, incrNum uint32) (newNum uint32, err error) {
	return m.IncrCounter(map[string]*dynamodb.AttributeValue{
		"name": {
			S: aws.String(name),
		},
	}, "count", incrNum)
}
