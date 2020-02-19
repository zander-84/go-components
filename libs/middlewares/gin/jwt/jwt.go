package CGinJwt

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"time"
)

var _jwt = new(Jwt)

type Jwt struct {
	conf           Conf
	AuthMiddleware *jwt.GinJWTMiddleware
	mutex          sync.Mutex
}

func New() *Jwt {
	return _jwt
}

func (j *Jwt) Init(user User, conf Conf) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	var err error
	j.conf = conf

	j.AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:            j.conf.Realm,
		SigningAlgorithm: j.conf.SigningAlgorithm,
		Key:              []byte(j.conf.Key),
		Timeout:          j.conf.Timeout,
		MaxRefresh:       j.conf.MaxRefresh,
		Authenticator:    user.Authenticator,
		Authorizator:     user.Authorizator,
		PayloadFunc:      user.PayloadFunc,
		Unauthorized:     user.Unauthorized,
		LoginResponse:    user.LoginResponse,
		LogoutResponse:   user.LogoutResponse,
		RefreshResponse:  user.RefreshResponse,
		IdentityHandler:  user.IdentityHandler,
		IdentityKey:      j.conf.IdentityKey,
		TokenLookup:      j.conf.TokenLookup,
		TokenHeadName:    j.conf.TokenHeadName,
		TimeFunc:         time.Now,

		HTTPStatusMessageFunc: nil,
		PrivKeyFile:           "",
		PubKeyFile:            "",
		SendCookie:            false,
		SecureCookie:          false,
		CookieHTTPOnly:        false,
		CookieDomain:          "",
		SendAuthorization:     false,
		DisabledAbort:         false,
		CookieName:            "",
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
}

func (j *Jwt) Middleware() gin.HandlerFunc {
	return j.AuthMiddleware.MiddlewareFunc()
}
