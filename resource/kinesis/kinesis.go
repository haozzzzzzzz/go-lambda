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

func GetKinesis(region string) (svc *kinesis.Kinesis, err error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if nil != err {
		logrus.Errorf("new aws session failed. \n%s.", err)
		return
	}

	// Create a Kinesis client with additional configuration
	svc = kinesis.New(sess)
	return
}
