package CAliOss

import "strings"

var _ossPrivate = new(OssPrivate)

type OssPrivate struct {
	Oss
}

func (this *OssPrivate) construct(conf Conf) *OssPrivate {
	this.conf = conf
	if this.conf.Dir != "" {
		this.conf.Dir = strings.Trim(this.conf.Dir, "/") + "/"
	}
	return this
}
func NewOssPrivate(conf Conf) *OssPrivate {
	_ossPrivate.construct(conf)
	return _ossPrivate
}
