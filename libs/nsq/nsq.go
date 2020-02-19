package CNsq

type Nsq struct {
	conf     Conf
	producer *Producer
	consumer *Consumer
}

//
func NewNsq(opts ...func(interface{})) *Nsq {
	var _nsq = new(Nsq)
	for _, opt := range opts {
		opt(_nsq)
	}
	_nsq.build()
	return _nsq
}
func BuildNsq(opts ...func(interface{})) interface{} {
	return NewNsq(opts...)
}
func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Nsq)
		g.conf = conf
		g.conf.SetDefault()
	}
}
func (v *Nsq) build() {
	v.producer = NewProducer()
	v.consumer = NewConsumer()
}

func (v *Nsq) Producer() *Producer {
	return v.producer
}

func (v *Nsq) Consumer() *Consumer {
	return v.consumer
}
