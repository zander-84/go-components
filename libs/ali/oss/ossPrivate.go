package CAliOss

var _ossPrivate = new(OssPrivate)

type OssPrivate struct {
	Oss
}

func NewOssPrivate(conf Conf) *OssPrivate {
	_ossPrivate.construct(conf)
	return _ossPrivate
}
