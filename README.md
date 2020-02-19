go components
===========================

## 组件通过配置文件更新
```go 
package main

import (
	"github.com/jinzhu/gorm"
    "github.com/zander-84/go-components"
	"time"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

)
func main(){
    c := C.NewComponents("./")
    c.Conf().Obj().(*viper.Viper).OnConfigChange(func(in fsnotify.Event) {
        c.Conf().ReloadBasicConf()
        c.Notify()
    })
}
```
