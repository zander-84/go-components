go components
===========================

## GinExample
```go
package main

import (
		"context"
    	"fmt"
    	"github.com/gin-gonic/gin"
    	"github.com/gin-gonic/gin/binding"
    	"github.com/zander-84/go-components"
    	"net/http"
    	"time"

)
type Login struct {
	User      string `form:"user" json:"user"  validate:"required,min=5"  comment:"用户名"`
	Passwords string `form:"password" json:"password" xml:"password" validate:"required,max=4" comment:"密码"`
}
func main(){
	c := C.NewComponents("./")
    gin.SetMode("debug")
  	r := gin.New()
  	binding.Validator = c.Validator()
    r.GET("/ping", func(c *gin.Context) {
		var loginData Login
		if err := c.ShouldBindQuery(&loginData); err != nil {
			fmt.Println(err.Error())  // {"password":"密码为必填字段","user":"用户名为必填字段"}
		} else {
			fmt.Println(loginData)
		}
		c.JSON(200, gin.H{
			"message": "pong",
			"data":    "success",
		})
	})
    r.Run() 
}
```
