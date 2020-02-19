package CGinCros

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"sync"
)

var _cors = new(Cors)

type Cors struct {
	conf Conf
	mutex sync.Mutex
}

func New() *Cors {
	return _cors
}

func (c *Cors) Init(conf Conf){
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.conf = conf
}

func (c *Cors) Middleware() gin.HandlerFunc {
	return cors.New(c.conf.Cors)
}
