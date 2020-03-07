package CHelper

import (
	"github.com/zander-84/go-components/libs/helper/file"
	"github.com/zander-84/go-components/libs/helper/format"
	"github.com/zander-84/go-components/libs/helper/request"
	"github.com/zander-84/go-components/libs/helper/security"
	"github.com/zander-84/go-components/libs/helper/time"
	CHelperTransform "github.com/zander-84/go-components/libs/helper/transform"
	"github.com/zander-84/go-components/libs/helper/type"
	"github.com/zander-84/go-components/libs/helper/unique"
	"sync"
)

type Helper struct {
	conf Conf
	once sync.Once

	unique    interface{}
	security  interface{}
	format    interface{}
	timeZone  interface{}
	file      interface{}
	slice     interface{}
	conv      interface{}
	httpCli   interface{}
	transform interface{}
}

func NewHelper(opts ...func(interface{})) *Helper {
	_helper := &Helper{}
	for _, opt := range opts {
		opt(_helper)
	}
	_helper.build()

	return _helper
}

func BulidHelper(opts ...func(interface{})) interface{} {
	return NewHelper(opts...)
}

func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Helper)
		g.conf = conf
		g.conf.SetDefault()
	}
}
func (v *Helper) build() {
	v.security = CHelperSecurity.NewSecurity()
	v.unique = CHelperUnique.NewUnique()
	v.format = CHelperFormat.NewFormat()
	v.timeZone = CHelperTime.NewTimeZone()
	v.file = CHelperFile.NewFile()
	v.slice = CHelperType.NewSlice()
	v.conv = CHelperType.NewConv()
	v.httpCli = CHelperRequest.NewHttpCli()
	v.transform = CHelperTransform.NewTransform()
}

func (v *Helper) Unique() *CHelperUnique.Unique {
	return v.unique.(*CHelperUnique.Unique)
}

func (v *Helper) Security() *CHelperSecurity.Security {
	return v.security.(*CHelperSecurity.Security)
}

func (v *Helper) Format() *CHelperFormat.Format {
	return v.format.(*CHelperFormat.Format)
}

func (v *Helper) TimeZone() *CHelperTime.TimeZone {
	return v.timeZone.(*CHelperTime.TimeZone)
}

func (v *Helper) File() *CHelperFile.File {
	return v.file.(*CHelperFile.File)
}

func (v *Helper) Slice() *CHelperType.Slice {
	return v.slice.(*CHelperType.Slice)
}

func (v *Helper) Conv() *CHelperType.Conv {
	return v.conv.(*CHelperType.Conv)
}

func (v *Helper) HttpCli() *CHelperRequest.HttpCli {
	return v.httpCli.(*CHelperRequest.HttpCli)
}

func (v *Helper) Transform() *CHelperTransform.Transform {
	return v.transform.(*CHelperTransform.Transform)
}
