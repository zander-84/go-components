

## Example
```go
package main

import (
	"github.com/zander-84/go-components"
	"time"
	"fmt"
	"github.com/zander-84/go-components/libs/cron"
	"github.com/nsqio/go-nsq"


)
type Login struct {
	User      string `form:"user" json:"user"  validate:"required,min=5"  comment:"用户名"`
	Passwords string `form:"password" json:"password" xml:"password" validate:"required,max=4" comment:"密码"`
}

func (l *Login) Run() {
	//fmt.Println(component.Get().Cron().Status())
	fmt.Println(l.User)
}
func main(){
	
	c := C.NewComponents("./")
   go func() {
   
   		for{
   			if  err := c.Nsq().Producer().MultiPublish(nil,c.Conf().Components.Nsq.ProducerTcpAddrs[0], "test2", [][]byte{{'h','e','l','l','o'},{'w','o','r','l','d'}}); err != nil {
   				fmt.Println("producer err ", err.Error())
   			} else {
   
   			}
   			time.Sleep(time.Second)
   		}
   
   	}()
   
   	if  err := c.Nsq().Consumer().ConsumeByLookupds(nil, "test2", "ch", c.Conf().Components.Nsq.ConsumerHttpAddrs, nsq.HandlerFunc(func(message *nsq.Message) error {
   			fmt.Printf("Got aa message: %s \n", message.Body)
   			message.Finish()
   			return nil
   		}), 2); err != nil {
   		fmt.Println("consumer err ", err.Error())
   	}
}
```
