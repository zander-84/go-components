package CAliOss

import "github.com/zander-84/go-components/libs/helper"

type Conf struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	Bucket          string
	Dir             string
	Host            string
	Helper          *CHelper.Helper
}
