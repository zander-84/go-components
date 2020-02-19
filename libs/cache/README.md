go components
===========================

## Example
```go
package main

import (
	"github.com/zander-84/go-components"
	"time"
	"fmt"

)
func main(){
	c := C.NewComponents("./")
    c.Cache().Set("name","zander",100*time.Second)
    var name string
    c.Cache().Get("name", &name)
    fmt.Println(name)

    var age int
    c.Cache().GetOrSet("age",&age, func() (value interface{}, err error) {
        return  18,nil
    },100*time.Second)
    fmt.Println(age)
}
```
