package kinesis

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/haozzzzzzzz/go-lambda/resource"
	"github.com/sirupsen/logrus"
)

func init() {
	resource.RegisterResource(resource.KinesisResourceType)
}

func NewSimpleKinesis(region string) (svc *kinesis.Kinesis, err error) {
	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if nil != err {
		logrus.Errorf("new aws session failed. %s.", err)
		return
	}

	svc = kinesis.New(ses)

	return
}

type KinesisClient struct {
	StreamName string
	Session    *session.Session
	Kinesis    *kinesis.Kinesis
}

func NewKinesis(sess *session.Session, streamName string) (kinesisClient *KinesisClient, err error) {
	svc := kinesis.New(sess)
	kinesisClient = &KinesisClient{
		StreamName: streamName,
		Session:    sess,
		Kinesis:    svc,
	}

	return
}
