package CNsq

import (
	"github.com/nsqio/go-nsq"
	"time"
)

var _producer = new(Producer)

type Producer struct{}

func NewProducer() *Producer {
	return _producer
}

func (this *Producer) Get(conf *nsq.Config, addr string) (producer *nsq.Producer, err error) {
	if conf == nil {
		conf = nsq.NewConfig()
	}

	if producer, err = nsq.NewProducer(addr, conf); err != nil {
		return nil, err
	}

	if err := producer.Ping(); err != nil {
		return nil, err
	}

	return producer, nil
}

func (this *Producer) Stop(producer *nsq.Producer) {
	producer.Stop()
}

func (this *Producer) Publish(producer *nsq.Producer, topic string, body []byte) error {
	return producer.Publish(topic, body)
}

func (this *Producer) MultiPublish(producer *nsq.Producer, topic string, body [][]byte) error {
	return producer.MultiPublish(topic, body)
}

//推动到延迟队列  适用于很多定时场景
func (this *Producer) DeferredPublish(producer *nsq.Producer, topic string, delay time.Duration, body []byte) error {
	return producer.DeferredPublish(topic, delay, body)
}
