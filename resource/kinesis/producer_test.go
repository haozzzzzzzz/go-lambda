package kinesis

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func TestKinesisProducer_PutRecord(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})
	if nil != err {
		t.Error(err)
		return
	}

	producer, err := NewKinesisProducer(sess, "KinesisTest")
	if nil != err {
		t.Error(err)
		return
	}

	producer.PutRecord([]byte("hello"))
	select {}
}
