package CGinSign

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/zander-84/go-components/libs/middlewares/gin/response"
	"io/ioutil"
	"sync"
)

var _sign = new(Sign)

type Sign struct {
	mutex sync.Mutex
	conf Conf
}

func New() *Sign {
	return _sign
}
func (r *Sign) Init(conf Conf) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	conf = conf.setDefault()

	r.conf = conf
}

func (s *Sign) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		if !s.checkSign(body, c.Request.Header.Get(s.conf.HeaderKey)) {
			CGinResponse.Resp.Forbidden(c)
		}

		if CGinResponse.HasResp(c) {
			CGinResponse.HandleResponse(c, false)
			c.Abort()
			return
		}else {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			c.Next()
		}

	}
}


func (s *Sign) checkSign(body []byte, headerSignVal string) bool {
	if len(headerSignVal) < 1 {
		return false
	}

	body = append(body, []byte(s.conf.Key)...)
	data := sha256.Sum256(body)

	sign := hex.EncodeToString(data[:])
	return sign == headerSignVal
}
