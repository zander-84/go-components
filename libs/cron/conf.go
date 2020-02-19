package CCron

type Conf struct {
	Location string
}

func (c *Conf) SetDefault() Conf{
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.Location == ""{
		c.Location = "Asia/Shanghai"
	}
}
