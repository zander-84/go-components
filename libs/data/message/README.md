go components
===========================

## 消息模型 Example
```go
package main

import (
    "fmt"
    "github.com/zander-84/go-components"
    "github.com/zander-84/go-components/libs/data/message"
    "github.com/zander-84/go-components/libs/data/queue"
    "math/rand"
    "time"
)
func main(){
    c := C.NewComponents("./")
	message := c.Message(CQueue.MaxPriorityQ)
  	go func() {
  		for i:=0;i<100;i++{
  			message.ProducePriority(fmt.Sprintf("1-%d",rand.Int()),1)
  			message.ProducePriority(fmt.Sprintf("2-%d",rand.Int()),2)
  			message.ProducePriority(fmt.Sprintf("3-%d",rand.Int()),3)
  			message.ProducePriority(fmt.Sprintf("4-%d",rand.Int()),4)
  
  			time.Sleep(time.Second)
  		}
  	}()
  	go func() {
  		message.Consume(func(mes *CMessage.Messages) error {
  			fmt.Println("consume1: ",mes.Data)
  			//mes.Finish()
  			panic("错误")
  			//time.Sleep(time.Second)
  			return nil
  		})
  	}()
  	go func() {
  		message.ConsumeFailQueue(func(mes *CMessage.Messages) error {
  			fmt.Println("failConsume3: ", mes.Data)
  			time.Sleep(time.Second)
  
  			return nil
  		})
  	}()
  
  	message.Consume(func(mes *CMessage.Messages) error {
  		fmt.Println("consume2: ", mes.Data)
  		return nil
  	})
}
```
