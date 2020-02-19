package CNsq

import "github.com/nsqio/go-nsq"

var _consumer = new(Consumer)

type Consumer struct {
}

func NewConsumer() *Consumer {
	return _consumer
}

func (n *Consumer) Get(conf *nsq.Config, topic string, channel string) (*nsq.Consumer, error) {

	if conf == nil {
		conf = nsq.NewConfig()
	}

	return nsq.NewConsumer(topic, channel, conf)
}

//  ConnectToNSQLookupds  循环 ConnectToNSQLookupd
func (n *Consumer) ConsumeByLookupds(conf *nsq.Config, topic string, channel string, address []string, handler nsq.Handler, concurrency int) error {
	var consumer *nsq.Consumer
	var err error
	if conf == nil {
		conf = nsq.NewConfig()
	}

	if consumer, err = nsq.NewConsumer(topic, channel, conf); err != nil {
		return err
	}
	consumer.AddConcurrentHandlers(handler, concurrency)
	return consumer.ConnectToNSQLookupds(address)
}

// ConnectToNSQDs 循环 ConnectToNSQD
func (n *Consumer) ConsumeByNSQDS(conf *nsq.Config, topic string, channel string, address []string, handler nsq.Handler, concurrency int)  error {
	var consumer *nsq.Consumer
	var err error
	if conf == nil {
		conf = nsq.NewConfig()
	}

	if consumer, err = nsq.NewConsumer(topic, channel, conf); err != nil {
		return err
	}
	consumer.AddConcurrentHandlers(handler, concurrency)
	return consumer.ConnectToNSQDs(address)


}

func (n *Consumer) MultiHandles(consumer *nsq.Consumer, handler nsq.Handler, concurrency int) {
	consumer.AddConcurrentHandlers(handler, concurrency)
}

func (n *Consumer) Stats(consumer *nsq.Consumer) *nsq.ConsumerStats {
	return consumer.Stats()
}
