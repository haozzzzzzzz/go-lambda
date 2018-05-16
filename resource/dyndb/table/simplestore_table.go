package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sirupsen/logrus"
)

// 简单存储模型
type SimpleStoreModel struct {
	PartitionKey string      `json:"partition_key"`
	SortKey      string      `json:"sort_key"`
	Value        interface{} `json:"value"`
}

type SimpleStoreTable struct {
	DynamoDBTable
}

func (m *SimpleStoreTable) AddSimpleStore(obj *SimpleStoreModel) (err error) {
	return m.PutItem(obj)
}

func (m *SimpleStoreTable) GetSimpleStore(partitionKey string, sortKey string) (obj *SimpleStoreModel, err error) {
	obj = &SimpleStoreModel{}
	err = m.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(m.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"partition_key": {
				S: aws.String(partitionKey),
			},
			"sort_key": {
				S: aws.String(sortKey),
			},
		},
	}, obj)
	if nil != err {
		logrus.Errorf("get simple store item failed. %s.", err)
		return
	}
	return
}

func (m *SimpleStoreTable) DeleteSimpleStore(partitionKey string, sortKey string) (err error) {
	_, err = m.Client.DeleteItemWithContext(m.Ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(m.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"partition_key": {
				S: aws.String(partitionKey),
			},
			"sort_key": {
				S: aws.String(sortKey),
			},
		},
	})
	if nil != err {
		logrus.Errorf("delete simple store table failed. %s.", err)
		return
	}
	return
}
