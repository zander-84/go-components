package CGinCommon

import "github.com/gin-gonic/gin"

var _common = new(Common)

type Common struct {
}

func New() *Common {
	return _common
}
func (c *Common) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("isGuest", true)
		c.Next()
	}
}
