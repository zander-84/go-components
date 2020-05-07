go components
===========================

## 定时任务 Example
```go
package main

import (
		"context"
    	"fmt"
    	"github.com/zander-84/go-components"
    	"net/http"
    	"os"
    	"time"
)


const TEXT1 = `甲方（签章）：                           乙方：中国xx股份有限公司XX分公司
甲方代表（签字）：                       乙方盖章：
签订日期：      年    月    日           签订日期：      年    月    日`

const TEXT2 = `aaaaaaaaaaaaaaaa`

func main(){
	c := C.NewComponents("./")
 	comp := NewCompound()
 	comp.Init(1200, 1200, "./src/simsun.ttf", 18, 72, 1.2)
 	if err := comp.AddTitle("移动业务靓号使用协议", 30); err != nil {
 		fmt.Println("add title error:", err.Error())
 	}
 
 	if err := comp.AddBody(TEXT1, 30, 1100); err != nil {
 		fmt.Println("add AddBody error:", err.Error())
 	}
 	if err := comp.AddImage("./a.png", 400, 400, 140, 0, false, false); err != nil {
 		fmt.Println("add img error:", err.Error())
 	}
 
 	if err := comp.AddImage("./a.png", 400, 400, 600, 0, false, true); err != nil {
 		fmt.Println("add img error:", err.Error())
 	}
 	comp.AddMarginSpace(30)
 
 	if err := comp.AddBody(TEXT2, 30, 1100); err != nil {
 		fmt.Println("add AddBody error:", err.Error())
 	}
 	if err := comp.AddImage("./src/a.png", 80, 20, 150, 95, true, false); err != nil {
 		fmt.Println("add img error:", err.Error())
 	}
 
 
 	buf := bytes.NewBufferString("")
 	//file, _ := os.Create("./src/dst2")
 	if err:=comp.Save(buf); err != nil {
 		fmt.Println("buf error:", err.Error())
 	}
}
```
