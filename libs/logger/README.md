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
	c := C.NewComponents("./")
    	c.Log().Debug("error")
    	c.Log().ExtendDebug(comLogger.Data{
    		TraceId:  "",
    		SpanId:   "",
    		Uid:      0,
    		Msg:      "",
    		Raw:      nil,
    		From:     "",
    		Duration: 0,
    	})
}
```
