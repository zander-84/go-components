package CAli

import (
	com_ali_oss "github.com/zander-84/go-components/libs/ali/oss"
	"github.com/zander-84/go-components/libs/helper"
)

type Conf struct {
	AccessKeyId     string //全局 AccessKeyId
	AccessKeySecret string //全局 AccessKeySecret
	Oss             com_ali_oss.Conf
	Helper          *CHelper.Helper
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	c.SetDefaultOss()
	return *c
}

func (c *Conf) SetDefaultBasic() {

}

func (c *Conf) SetDefaultOss() {
	if c.Oss.AccessKeyId == "" && c.Oss.AccessKeySecret == "" {
		c.Oss.AccessKeyId = c.AccessKeyId
		c.Oss.AccessKeySecret = c.AccessKeySecret
	}
	c.Oss.Helper = c.Helper
}
