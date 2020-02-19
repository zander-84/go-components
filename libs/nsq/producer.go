package CNsq

import (
	"github.com/nsqio/go-nsq"
	"time"
)

var _producer = new(Producer)

type Producer struct {
}

func NewProducer() *Producer {
	return _producer
}

func (n *Producer) Get(conf *nsq.Config, addr string, ) (*nsq.Producer, error) {
	var producer *nsq.Producer
	var err error

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

func (n *Producer) Publish(conf *nsq.Config, addr string, topic string, body []byte) error {
	var producer *nsq.Producer
	var err error
	if producer, err = n.Get(conf, addr); err != nil {
		return err
	}
	defer producer.Stop()
	return producer.Publish(topic, body)
}

func (n *Producer) MultiPublish(conf *nsq.Config, addr string, topic string, body [][]byte) error {
	var producer *nsq.Producer
	var err error
	if producer, err = n.Get(conf, addr); err != nil {
		return err
	}
	defer producer.Stop()
	return producer.MultiPublish(topic, body)
}

//推动到延迟队列  适用于很多定时场景
func (n *Producer) DeferredPublish(conf *nsq.Config, addr string, topic string, delay time.Duration, body []byte) error {
	var producer *nsq.Producer
	var err error
	if producer, err = n.Get(conf, addr); err != nil {
		return err
	}
	defer producer.Stop()
	return producer.DeferredPublish(topic, delay, body)
}
