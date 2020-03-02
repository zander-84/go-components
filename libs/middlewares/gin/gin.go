package CGin

import (
	"github.com/gin-gonic/gin"
	"github.com/zander-84/go-components/libs/middlewares/gin/common"
	"github.com/zander-84/go-components/libs/middlewares/gin/cors"
	"github.com/zander-84/go-components/libs/middlewares/gin/jwt"
	"github.com/zander-84/go-components/libs/middlewares/gin/rate-limit"
	"github.com/zander-84/go-components/libs/middlewares/gin/response"
	"github.com/zander-84/go-components/libs/middlewares/gin/sign"
	"sync"
)

type GinMiddleWare struct {
	once sync.Once
}

func New() *GinMiddleWare {
	return new(GinMiddleWare)
}

//
//______________________________________________________________________
func (c *GinMiddleWare) InitCors(conf CGinCros.Conf) {
	CGinCros.New().Init(conf)
}

func (c *GinMiddleWare) Cors() gin.HandlerFunc {
	return CGinCros.New().Middleware()
}

//
//______________________________________________________________________
func (c *GinMiddleWare) InitRatelimiter(conf CGinRateLimit.Conf) {
	CGinRateLimit.New().Init(conf)
}
func (c *GinMiddleWare) Ratelimiter() gin.HandlerFunc {
	return CGinRateLimit.New().Middleware()
}

//
//______________________________________________________________________
func (c *GinMiddleWare) InitResponse(debug bool) {
	CGinResponse.MiddleResp.Init(debug)
}

func (c *GinMiddleWare) Response() gin.HandlerFunc {
	return CGinResponse.MiddleResp.Middleware()
}

//
//______________________________________________________________________
func (c *GinMiddleWare) Common() gin.HandlerFunc {
	return CGinCommon.New().Middleware()
}

//
//______________________________________________________________________
func (c *GinMiddleWare) InitJwt(user CGinJwt.User, conf CGinJwt.Conf) *CGinJwt.Jwt {
	return CGinJwt.New().Init(user, conf)
}

func (c *GinMiddleWare) Jwt(j *CGinJwt.Jwt) gin.HandlerFunc {
	return j.Middleware()
}

//
//______________________________________________________________________
func (c *GinMiddleWare) InitSign(conf CGinSign.Conf) {
	CGinSign.New().Init(conf)
}

func (c *GinMiddleWare) Sign() gin.HandlerFunc {
	return CGinSign.New().Middleware()
}
