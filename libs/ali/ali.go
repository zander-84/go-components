package CAli

import (
	oss2 "github.com/zander-84/go-components/libs/ali/oss"
	"github.com/zander-84/go-components/libs/helper"
	"sync"
)

type Ali struct {
	oss        *oss2.Oss
	ossPrivate *oss2.OssPrivate

	conf Conf
	once sync.Once
}

//
func NewAli(opts ...func(interface{})) *Ali {
	var _ali = &Ali{}
	for _, opt := range opts {
		opt(_ali)
	}
	_ali.build()
	return _ali
}

func BuildAli(opts ...func(interface{})) interface{} {
	return NewAli(opts...)
}

func SetConfig(conf Conf, helper *CHelper.Helper) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Ali)
		g.conf = conf
		g.conf.Helper = helper
		g.conf.SetDefault()
	}
}

func (this *Ali) construct(conf Conf, helper *CHelper.Helper) *Ali {
	this.conf = conf
	this.conf.Helper = helper
	this.conf.SetDefault()
	this.build()
	return this
}

func (this *Ali) build() {
	this.oss = oss2.New(this.conf.Oss)
	this.ossPrivate = oss2.NewOssPrivate(this.conf.OssPrivate)
}

func (this *Ali) Oss() *oss2.Oss {
	return this.oss
}

func (this *Ali) OssPrivate() *oss2.OssPrivate {
	return this.ossPrivate
}
