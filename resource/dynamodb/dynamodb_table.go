package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

type DynamoDBTable struct {
	TableName string
	Client    *dynamodb.DynamoDB
}

func (m *DynamoDBTable) GetItem(input *dynamodb.GetItemInput, item interface{}) (err error) {
	input.TableName = aws.String(m.TableName)
	getOutput, err := m.Client.GetItem(input)
	if nil != err {
		logrus.Errorf("get output failed. %s.", err)
		return
	}

	err = dynamodbattribute.UnmarshalMap(getOutput.Item, item)
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

	_, err = m.Client.PutItem(&dynamodb.PutItemInput{
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
	output, err := m.Client.Query(input)
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

func (m *DynamoDBTable) IncrCounter(key map[string]*dynamodb.AttributeValue, fieldName string, incrNum uint32) (err error) {
	output, err := m.Client.UpdateItem(&dynamodb.UpdateItemInput{
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

	fmt.Println(output)
	return
}
