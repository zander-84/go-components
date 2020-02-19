package CGinCros

import "github.com/gin-contrib/cors"

type Conf struct {
	Cors cors.Config
}

func (c Conf) Default() Conf {
	conf :=Conf{}
	conf.Cors = cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Accept", "Accept-Language", "Authorization", "X-Rate-Limit-Limit", "X-Rate-Limit-Duration", "X-Rate-Limit-Request-Forwarded-For", "X-Rate-Limit-Request-Remote-Addr"},
		ExposeHeaders:    []string{"X-My-Header"},
		AllowCredentials: false,
	}

	return conf
}
