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
func (this *Consumer) Stop(consumer *nsq.Consumer) {
	consumer.Stop()
}

//  ConnectToNSQLookupds  循环 ConnectToNSQLookupd
func (n *Consumer) ConsumeByLookupds(consumer *nsq.Consumer, address []string, handler nsq.Handler, concurrency int) error {
	consumer.AddConcurrentHandlers(handler, concurrency)
	return consumer.ConnectToNSQLookupds(address)
}

// ConnectToNSQDs 循环 ConnectToNSQD
func (n *Consumer) ConsumeByNSQDS(consumer *nsq.Consumer, address []string, handler nsq.Handler, concurrency int) error {
	consumer.AddConcurrentHandlers(handler, concurrency)
	return consumer.ConnectToNSQDs(address)
}

//
func (n *Consumer) ConsumeByNSQD(consumer *nsq.Consumer, address string, handler nsq.Handler, concurrency int) error {
	consumer.AddConcurrentHandlers(handler, concurrency)
	return consumer.ConnectToNSQD(address)
}

func (n *Consumer) MultiHandles(consumer *nsq.Consumer, handler nsq.Handler, concurrency int) {
	consumer.AddConcurrentHandlers(handler, concurrency)
}

func (n *Consumer) Stats(consumer *nsq.Consumer) *nsq.ConsumerStats {
	return consumer.Stats()
}
