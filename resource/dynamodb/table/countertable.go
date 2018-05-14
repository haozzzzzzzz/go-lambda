package table

import (
	"github.com/aws/aws-sdk-go/aws"
	dynamodb2 "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/haozzzzzzzz/go-lambda/resource/dynamodb"
)

// 计数器
type CounterModel struct {
	Name  string `json:"name"`
	Count uint32 `json:"count"`
}

type CounterTable struct {
	dynamodb.DynamoDBTable
}

func (m *CounterTable) Incr(name string, incrNum uint32) (newNum uint32, err error) {
	return m.IncrCounter(map[string]*dynamodb2.AttributeValue{
		"name": {
			S: aws.String(name),
		},
	}, "count", incrNum)
}
