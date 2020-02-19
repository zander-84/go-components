package CGinResponse

import "github.com/gin-gonic/gin"

const Key  = "resp"
var (
	Success = "成功"
	SystemSpaceErr = "系统空间错误"
	UserSpaceErr = "用户空间错误"
	UnauthorizenErr = "未登入"
	RateLimitErr = "访问过于频繁"
	FrobiddenErr = "禁止访问"
	SignErr = "数据错误"
)


type Response interface {
	//成功
	Success(c *gin.Context, data interface{}, debugs ...interface{})

	// 系统空间错误
	SystemSpaceErr(c *gin.Context, debug interface{})

	// 用户空间错误
	UserSpaceErr(c *gin.Context, data interface{}, debugs ...interface{})

	// 未认证
	Unauthorized(c *gin.Context)

	// 限速
	RateLimit(c *gin.Context)

	// 禁止
	Forbidden(c *gin.Context)

	// 自定义
	Custom(c *gin.Context, code int, msg interface{}, data interface{}, debugs ...interface{})

	Middleware() gin.HandlerFunc

	hasResp(c *gin.Context) bool

	Init(bool2 bool)
}
