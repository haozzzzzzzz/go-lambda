package table

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

type DynamoDBTableError struct {
	Message string `json:"message"`
}

func (m *DynamoDBTableError) Error() string {
	return m.Message
}

var NilItem = &DynamoDBTableError{
	Message: "nil table item",
}

type DynamoDBTable struct {
	TableName string
	Ctx       aws.Context
	Client    *dynamodb.DynamoDB
}

func (m *DynamoDBTable) GetItem(input *dynamodb.GetItemInput, item interface{}) (err error) {
	input.TableName = aws.String(m.TableName)
	output, err := m.Client.GetItemWithContext(m.Ctx, input)
	if nil != err {
		logrus.Errorf("get output failed. %s.", err)
		return
	}

	// 没有记录
	if len(output.Item) == 0 {
		err = NilItem
		return
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, item)
	if nil != err {
		logrus.Errorf("unmarshal item failed. %s.", err)
		return
	}

	return
}

func (m *DynamoDBTable) PutItem(item interface{}) (err error) {
	attributeValue, err := dynamodbattribute.MarshalMap(item)
	if nil != err {
		logrus.Errorf("marshal item failed. %s.", err)
		return
	}

	_, err = m.Client.PutItemWithContext(m.Ctx, &dynamodb.PutItemInput{
		TableName: aws.String(m.TableName),
		Item:      attributeValue,
	})
	if nil != err {
		logrus.Errorf("put item failed. %s.", err)
		return
	}

	return
}

func (m *DynamoDBTable) Query(input *dynamodb.QueryInput, records interface{}) (err error) {
	input.TableName = aws.String(m.TableName)
	output, err := m.Client.QueryWithContext(m.Ctx, input)
	if nil != err {
		logrus.Errorf("query items failed. %s.", err)
		return
	}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, records)
	if nil != err {
		logrus.Errorf("unmarshal items failed. %s.", err)
		return
	}
	return
}

// 增加计数
func (m *DynamoDBTable) IncrCounter(key map[string]*dynamodb.AttributeValue, fieldName string, incrNum uint32) (newNum uint32, err error) {
	output, err := m.Client.UpdateItemWithContext(m.Ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(m.TableName),
		Key:              key,
		UpdateExpression: aws.String("ADD #field :incr"),
		ExpressionAttributeNames: map[string]*string{
			"#field": aws.String(fieldName),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": {
				N: aws.String(fmt.Sprintf("%d", incrNum)),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	})
	if nil != err {
		logrus.Errorf("update counter failed. %s.", err)
		return
	}

	attrValue := output.Attributes[fieldName]
	newNum64, err := strconv.ParseInt(*attrValue.N, 10, 32)
	if nil != err {
		logrus.Errorf("parse attribute value string to int failed. %s.", err)
		return
	}

	newNum = uint32(newNum64)
	return
}
