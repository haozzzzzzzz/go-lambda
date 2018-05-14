package table

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	dynamodb2 "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/haozzzzzzzz/go-lambda/resource/dynamodb"
	"github.com/sirupsen/logrus"
)

// 首页缓存
type CacheModel struct {
	Key      string `json:"key"`
	Record   []byte `json:"record"`
	ExpireAt int64  `json:"expire_at"`
}

type CacheTable struct {
	dynamodb.DynamoDBTable
}

func (m *CacheTable) SetNotExist(key string, obj interface{}, ttl time.Duration) (success bool, err error) {

	return
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
