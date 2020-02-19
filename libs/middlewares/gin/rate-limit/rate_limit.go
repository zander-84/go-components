package CGinRateLimit

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/zander-84/go-components/libs/middlewares/gin/response"
	"sync"
)

var _rateLimiter = new(RateLimiter)

type RateLimiter struct {
	conf  Conf
	mutex sync.Mutex
}

func New() *RateLimiter {
	return _rateLimiter
}

var lmt *limiter.Limiter

func (r *RateLimiter) Init(conf Conf) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.conf = conf

	if r.conf.Max < 1 {
		r.conf.Max = 10
	}

	lmt = tollbooth.NewLimiter(r.conf.Max, nil)
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	lmt.SetBurst(r.conf.Burst)
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)

		if r.conf.RemoveHeader {
			c.Writer.Header().Del("X-Rate-Limit-Request-Forwarded-For")
			c.Writer.Header().Del("X-Rate-Limit-Request-Remote-Addr")
			c.Writer.Header().Del("X-Rate-Limit-Duration")
			c.Writer.Header().Del("X-Rate-Limit-Limit")
		}

		if httpError != nil {
			CGinResponse.Resp.RateLimit(c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
