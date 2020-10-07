package CAli

import (
	com_ali_oss "github.com/zander-84/go-components/libs/ali/oss"
	"github.com/zander-84/go-components/libs/helper"
)

type Conf struct {
	AccessKeyId     string //全局 AccessKeyId
	AccessKeySecret string //全局 AccessKeySecret
	Oss             com_ali_oss.Conf
	OssPrivate      com_ali_oss.Conf
	Helper          *CHelper.Helper
}

func (this *Conf) SetDefault() Conf {
	this.SetDefaultBasic()
	this.SetDefaultOss()
	return *this
}

func (this *Conf) SetDefaultBasic() {

}

func (this *Conf) SetDefaultOss() {
	if this.Oss.AccessKeyId == "" && this.Oss.AccessKeySecret == "" {
		this.Oss.AccessKeyId = this.AccessKeyId
		this.Oss.AccessKeySecret = this.AccessKeySecret
	}
	if this.OssPrivate.AccessKeyId == "" && this.OssPrivate.AccessKeySecret == "" {
		this.OssPrivate.AccessKeyId = this.AccessKeyId
		this.OssPrivate.AccessKeySecret = this.AccessKeySecret
	}
	this.Oss.Helper = this.Helper
	this.OssPrivate.Helper = this.Helper
}
