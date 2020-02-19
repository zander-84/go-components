package CMemory

type Conf struct {
	Expiration        int //分钟
	CleanupInterval int //分钟
}

func (c *Conf) SetDefault() Conf{
	c.SetDefaultBasic()
	return *c
}

func (c *Conf) SetDefaultBasic() {

}
