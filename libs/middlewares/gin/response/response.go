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

func (r *defaultResponse) SuccessAction(c *gin.Context) {
	c.Set(Key, &defaultResponse{
		Code: http.StatusOK,
		Msg:  Codes[http.StatusOK],
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
