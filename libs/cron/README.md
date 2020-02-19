go components
===========================

## 定时任务 Example
```go
package main

import (
		"context"
    	"fmt"
    	"github.com/zander-84/go-components"
    	"github.com/zander-84/go-components/libs/cron"
    	"net/http"
    	"os"
    	"time"
)

type Login struct {
	User      string `form:"user" json:"user"  validate:"required,min=5"  comment:"用户名"`
	Passwords string `form:"password" json:"password" xml:"password" validate:"required,max=4" comment:"密码"`
}

func (l *Login) Run() {
	fmt.Println(C.Get().Cron().Status())
	fmt.Println(l.User, l.Passwords)
}

func main(){
    	c := C.NewComponents("./")
    	cron := c.Cron()
    	if err := cron.AddJob(&CCron.Job{
    		ID:   "test1",
    		Desc: "测试1",
    		Spec: "* * * * * *",
    		Cmd: &Login{
    			User:      "zander",
    			Passwords: "pawd",
    		},
    		Obj: nil,
    	}); err != nil {
    		fmt.Println(err)
    	}
    	if err := cron.AddJob(&CCron.Job{
    		ID:   "test2",
    		Desc: "测试1",
    		Spec: "* * * * * *",
    		Cmd: CCron.CmdFunc(func() {
    			fmt.Println("hello word")
    		}),
    		Obj: nil,
    	}); err != nil {
    		fmt.Println(err)
    	}
    	cron.Start()
    	cron.Restart("test2")
    	cron.Remove("test1")
    	select {}
}
```
