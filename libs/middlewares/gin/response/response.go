package CGinResponse

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var Resp Response = &defaultResponse{}
var conf Conf = Conf{Debug: false}
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
		Msg:   Codes[http.StatusOK],
		Data:  data,
		Debug: debugs,
	})
}

func (r *defaultResponse) SystemSpaceErr(c *gin.Context, debug interface{}) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusInternalServerError,
		Msg:   Codes[http.StatusInternalServerError],
		Data:  nil,
		Debug: debug,
	})
}

func (r *defaultResponse) UserSpaceErr(c *gin.Context, data interface{}, debugs ...interface{}) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusBadRequest,
		Msg:   Codes[http.StatusBadRequest],
		Data:  data,
		Debug: debugs,
	})
}

func (r *defaultResponse) ParamsErr(c *gin.Context, data interface{}, debugs ...interface{}) {
	c.Set(Key, &defaultResponse{
		Code:  StatusParamsErr,
		Msg:   Codes[StatusParamsErr],
		Data:  data,
		Debug: debugs,
	})
}

func (r *defaultResponse) RateLimit(c *gin.Context) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusTooManyRequests,
		Msg:   Codes[http.StatusTooManyRequests],
		Data:  nil,
		Debug: nil,
	})
}

func (r *defaultResponse) Forbidden(c *gin.Context) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusForbidden,
		Msg:   Codes[http.StatusForbidden],
		Data:  nil,
		Debug: nil,
	})
}

func (r *defaultResponse) Unauthorized(c *gin.Context) {
	c.Set(Key, &defaultResponse{
		Code:  http.StatusUnauthorized,
		Msg:   Codes[http.StatusUnauthorized],
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
			if rec := recover(); rec != nil {
				r.SystemSpaceErr(c, rec)
				HandleResponse(c, conf.Debug)
				c.Abort()
				return
			}
		}()

		if HasResp(c) {
			HandleResponse(c, conf.Debug)
			c.Abort()
			return
		} else {
			c.Next()
			if HasResp(c) {
				HandleResponse(c, conf.Debug)
				c.Abort()
				return
			}
		}
	}
}

func HasResp(c *gin.Context) bool {
	if _, ok := c.Get(Key); ok {
		return true
	}
	return false
}

func HandleResponse(c *gin.Context, debug bool) {
	if data, ok := c.Get(Key); ok {
		if data2, ok2 := data.(*defaultResponse); ok2 {
			body := gin.H{
				"code": data2.Code,
				"msg":  data2.Msg,
				"data": data2.Data,
			}
			if debug {
				body["debug"] = data2.Debug
				c.IndentedJSON(200, body)
			} else {
				c.JSON(200, body)
			}
		}
	}
}
