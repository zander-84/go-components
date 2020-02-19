package CGinResponse

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var Resp Response = &defaultResponse{}
var conf Conf = Conf{Debug:false}
var mutex sync.Mutex

type defaultResponse struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Data  interface{} `json:"data"`
	Debug interface{} `json:"debug"`
}

func (r *defaultResponse) Success(c *gin.Context, data interface{}, debugs ...interface{}) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusOK,
		Msg:   Success,
		Data:  data,
		Debug: debugs,
	})
}

func (r *defaultResponse) SystemSpaceErr(c *gin.Context, debug interface{}) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusInternalServerError,
		Msg:   SystemSpaceErr,
		Data:  nil,
		Debug: debug,
	})
}

func (r *defaultResponse) UserSpaceErr(c *gin.Context, data interface{}, debugs ...interface{}) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusBadRequest,
		Msg:   UserSpaceErr,
		Data:  data,
		Debug: debugs,
	})
}

func (r *defaultResponse) RateLimit(c *gin.Context) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusTooManyRequests,
		Msg:   RateLimitErr,
		Data:  nil,
		Debug: nil,
	})
}

func (r *defaultResponse) Forbidden(c *gin.Context) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusForbidden,
		Msg:   FrobiddenErr,
		Data:  nil,
		Debug: nil,
	})
}

func (r *defaultResponse) Unauthorized(c *gin.Context) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusUnauthorized,
		Msg:   UnauthorizenErr,
		Data:  nil,
		Debug: nil,
	})
}

func (r *defaultResponse) Custom(c *gin.Context, code int, msg interface{}, data interface{}, debugs ...interface{}) {
	c.Set(Key, &defaultResponse{
		Code:  code,
		Msg:   msg,
		Data:  data,
		Debug: debugs,
	})
}

func (r *defaultResponse) Init(debug bool) {
	mutex.Lock()
	defer mutex.Unlock()

	conf.Debug = debug
}

func (r *defaultResponse) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//____ 接受panic
		defer func() {
			if rec := recover(); r != nil {
				r.SystemSpaceErr(c, rec)
				r.handleResponse(c, conf.Debug)
			}
		}()

		if r.hasResp(c) {
			r.handleResponse(c, conf.Debug)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func (r *defaultResponse) handleResponse(c *gin.Context, debug bool) {
	if data, ok := c.Get(Key); ok {
		if data2, ok2 := data.(*defaultResponse); ok2 {
			body := gin.H{
				"code": data2.Code,
				"msg":  data2.Msg,
				"data": data2.Data,
			}
			if debug {
				body["debug"] = data2.Debug
			}

			c.JSON(data2.Code, body)
		}
	}
}
func (r *defaultResponse) hasResp(c *gin.Context) bool {
	if _, ok := c.Get(Key); ok {
		return true
	}
	return false
}
