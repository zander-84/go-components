package CGinSign

type Conf struct {
	Key string
	HeaderKey string
}

func (c Conf) setDefault() Conf {
	if c.Key == ""{
		c.Key = "hello gin"
	}

	if c.HeaderKey == ""{
		c.HeaderKey = "sign"
	}

	return c
}
