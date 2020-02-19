go components
===========================

#### 线程测试需要注意不要开太多优先级 Example
```go
package main

import (
    "context"
    "fmt"
    "github.com/zander-84/go-components"
    "github.com/zander-84/go-components/libs/worker"
    "net/http"
    "os"
    "time"
)

type Login struct {
	User      string `form:"user" json:"user"  validate:"required,min=5"  comment:"用户名"`
	Passwords string `form:"password" json:"password" xml:"password" validate:"required,max=4" comment:"密码"`
}

func (l *Login) Run() error{
	fmt.Println(l.User, l.Passwords)
    return nil
}

func main(){
    	c := C.NewComponents("./")
        	fmt.Println("start.............")
        
        	for i := 0; i < 10000000; i++ {
        		tmp:=i
        		c.Worker().AddJob(CWorker.JobFunc(func() error {
        			fmt.Printf("%s %d \n","hello word",tmp)
        			return nil
        		}), 1)
        	}
            c.Worker().AddJob(&Login{
					User:      fmt.Sprintf("%s %d","zander",1),
					Passwords: fmt.Sprintf("%s %d","pawd",1),
				},1)

    	select {}
}
```
