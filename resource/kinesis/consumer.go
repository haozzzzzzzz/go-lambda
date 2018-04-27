package kinesis

import (
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type KinesisConsumer struct {
	KinesisClient
	ConsumerGroup              string
	ConsumerChannel            chan *kinesis.Record
	ConsumerChannelBufferCount uint32
}

// consumer
func (m *KinesisConsumer) RunConsumer() (err error) {
	if m.ConsumerChannel != nil {
		return
	}

	m.ConsumerChannel = make(chan *kinesis.Record, m.ConsumerChannelBufferCount)

	return
}

func (m *KinesisConsumer) GetRecord() chan *kinesis.Record {
	return m.ConsumerChannel
}
