package CGinResponse

import "github.com/gin-gonic/gin"

var MiddleResp = &middleResponse{}

type middleResponse struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Data  interface{} `json:"data"`
	Debug interface{} `json:"debug"`
}

func (r *middleResponse) Init(debug bool) {
	mutex.Lock()
	defer mutex.Unlock()

	conf.Debug = debug
}

func (r *middleResponse) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//____ 接受panic
		defer func() {
			if rec := recover(); rec != nil {
				Resp.SystemSpaceErr(c, rec)
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
