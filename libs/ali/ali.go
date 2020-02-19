package CAli

import (
	oss2 "github.com/zander-84/go-components/libs/ali/oss"
	"github.com/zander-84/go-components/libs/helper"
	"sync"
)



type Ali struct {
	oss *oss2.Oss

	conf Conf
	once sync.Once
}

//
func NewAli(opts ...func(interface{}))*Ali {
	var _ali = &Ali{}
	for _,opt := range opts{
		opt(_ali)
	}
	_ali.build()
	return _ali
}

func BuildAli(opts ...func(interface{})) interface{} {
	return NewAli(opts...)
}

func SetConfig(conf Conf,helper *CHelper.Helper) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Ali)
		g.conf = conf
		g.conf.Helper = helper
		g.conf.SetDefault()
	}
}


func (v *Ali) construct(conf Conf, helper *CHelper.Helper) *Ali {
	v.conf = conf
	v.conf.Helper = helper
	v.conf.SetDefault()
	v.build()
	return v
}

func (v *Ali) build() {
	v.oss = oss2.New(v.conf.Oss)
}

func (v *Ali) Oss() *oss2.Oss {
	return v.oss
}
