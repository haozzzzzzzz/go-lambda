package table

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	dynamodb2 "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/haozzzzzzzz/go-lambda/resource/dynamodb"
	"github.com/sirupsen/logrus"
)

// 首页缓存
type CacheModel struct {
	Key      string      `json:"key"`
	Record   interface{} `json:"record"`
	ExpireAt int64       `json:"expire_at"`
}

type CacheTable struct {
	dynamodb.DynamoDBTable
}

func (m *CacheTable) SetNx(key string, obj interface{}, expireAt int64) (success bool, err error) {
	attributeValue, err := dynamodbattribute.MarshalMap(&CacheModel{
		Key:      key,
		Record:   obj,
		ExpireAt: expireAt,
	})
	output, err := m.Client.PutItemWithContext(m.Ctx, &dynamodb2.PutItemInput{
		TableName:           aws.String(m.TableName),
		Item:                attributeValue,
		ConditionExpression: aws.String("attribute_not_exists(key)"), // 做到这里
	})
	if nil != err {
		logrus.Errorf("put cache item failed. %s.", err)
		return
	}

	fmt.Println(output)

	return
}

func (m *CacheTable) SetNxTTL(key string, obj interface{}, ttl time.Duration) (bool, error) {
	return m.SetNx(key, obj, time.Now().Add(ttl).Unix())
}

func (m *CacheTable) Get(key string, obj interface{}) (err error) {
	err = m.GetItem(&dynamodb2.GetItemInput{
		Key: map[string]*dynamodb2.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
	}, obj)

	if nil != err {
		logrus.Errorf("get cache item failed. %s.", err)
		return
	}

	return
}

func (m *CacheTable) Delete(key string) (err error) {
	_, err = m.Client.DeleteItem(&dynamodb2.DeleteItemInput{
		TableName: aws.String(m.TableName),
		Key: map[string]*dynamodb2.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
	})
	if nil != err {
		logrus.Errorf("delete cache item failed. %s.", err)
		return
	}
	return
}
