package CCron

type Conf struct {
	TimeZone string
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.TimeZone == "" {
		c.TimeZone = "Asia/Shanghai"
	}
}
