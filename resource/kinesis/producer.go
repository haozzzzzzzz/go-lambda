package kinesis

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/haozzzzzzzz/go-rapid-development/utils/id"
	"github.com/sirupsen/logrus"
)

type KinesisProducer struct {
	KinesisClient
	ProducerChannel            chan *kinesis.PutRecordInput
	ProducerChannelBufferCount uint32
}

func NewKinesisProducer(sess *session.Session, streamName string) (producer *KinesisProducer, err error) {
	client, err := NewKinesis(sess, streamName)
	if nil != err {
		logrus.Errorf("new kinesis client failed. %s", err)
		return
	}

	producer = &KinesisProducer{
		KinesisClient:              *client,
		ProducerChannelBufferCount: 3000,
	}
	return
}

// producer
func (m *KinesisProducer) RunProducer() {
	if m.ProducerChannel != nil {
		return
	}

	m.ProducerChannel = make(chan *kinesis.PutRecordInput, m.ProducerChannelBufferCount)
	go func() {
		for {
			fmt.Println("Go routine running...")
			select {
			case input := <-m.ProducerChannel:
				// 启用协程进行处理，必要的时候，对协程数量进行控制
				go func(innerInput *kinesis.PutRecordInput) {
					_, err := m.Kinesis.PutRecord(innerInput)
					if err != nil {
						logrus.Errorf("kinesis put record failed. %s", err)
					}
				}(input)
			}
		}
	}()
}

func (m *KinesisProducer) PutRecord(data []byte) {
	partitionKey := id.UniqueID()
	input := &kinesis.PutRecordInput{
		Data:         data,
		StreamName:   &m.StreamName,
		PartitionKey: &partitionKey,
	}
	m.ProducerChannel <- input
}
