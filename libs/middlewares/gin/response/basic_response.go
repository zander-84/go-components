package CGinResponse

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const Key  = "resp"

var Codes = map[int]string{
	http.StatusOK : "成功",
	http.StatusInternalServerError : "系统空间错误",
	http.StatusBadRequest : "用户空间错误",
	http.StatusTooManyRequests : "访问过于频繁",
	http.StatusForbidden: "禁止访问",
	http.StatusUnauthorized: "未登入",

	// 自定义错误码
}



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

	Init(bool2 bool)
}
