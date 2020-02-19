go components
===========================

## Example
```go
package main

import (
	"github.com/zander-84/go-components"
	comLogger "github.com/zander-84/go-components/libs/logger"
	"time"
	"fmt"

)
func main(){
   	com := component.New("./")
   	com.MemoryCache().Set("name","zander",100*time.Second)
   	var name string
   	com.MemoryCache().Get("name", &name)
   	fmt.Println(name)
   
   	var age int
   	com.MemoryCache().GetOrSet("age",&age, func() (value interface{}, err error) {
   		return  18,nil
   	},100*time.Second)
   	fmt.Println(age)
}
```
