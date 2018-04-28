package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/haozzzzzzzz/go-lambda/resource"
	"github.com/sirupsen/logrus"
)

func init() {
	resource.RegisterResource(resource.DynamoDBResourceType)
}

func GetDynamodb(region string) (db *dynamodb.DynamoDB, err error) {
	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if nil != err {
		logrus.Errorf("new aws session failed. \n%s.", ses)
		return
	}

	db = dynamodb.New(ses)
	return
}
