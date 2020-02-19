package CNsq

type Conf struct {
	ProducerTcpAddrs  []string //
	ProducerHttpAddrs []string //
	ConsumerTcpAddrs  []string //
	ConsumerHttpAddrs []string // Lookupd
}

func (c *Conf) SetDefault() Conf{
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {

}
