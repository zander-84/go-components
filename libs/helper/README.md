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
    fmt.Println(com.Helper().Unique().ID())
}
```
