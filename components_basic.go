package C

import (
	"fmt"
	"strings"
	"sync"
)

type BasicComponents interface {
	// 添加
	Attach(name string, data interface{}) error

	// 移除
	Detach(name string)

	// 是否存在 yes:存在  no: 不存在
	IsExist(name string) bool

	// 获取组件
	Get(name string) interface{}

	// 打印组件
	Print()

	// 通知
	// 	com.Conf().Obj().(*viper.Viper).OnConfigChange(func(in fsnotify.Event) {
	//		com.Conf().ReloadBasicConf()
	//		com.Notify()
	//	})
	Notify(keys ...string)
}

var _ BasicComponents = new(c)

type c struct {
	objs  map[string]interface{} // 组件
	mutex sync.Mutex
	container Container              // 构建容器
}

func (this *c) Attach(name string, data interface{}) error {
	if this.IsExist(name) || this.isInternalKey(name) {
		return ErrKeyExist
	}
	this.attach(name, data)
	return nil
}

// 判断健
func (this *c) isInternalKey(name string) bool {
	return strings.HasPrefix(name, "inner.")
}

func (this *c) attach(name string, data interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.objs[name] = data
}

func (this *c) Detach(name string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.objs, name)
	this.container.delCache(name)
}

func (this *c) Get(name string) interface{} {
	return this.objs[name]
}

func (this *c) IsExist(name string) bool {
	_, ok := this.objs[name]
	return ok
}

// 通知。
// 组件采用懒加载模式，当配置文件更新时候只需要移除组件,再次使用自动更新组件
// 配合使用 conf.obj.OnConfigChange
func (this *c) Notify(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if key==fieldWorker{
				continue
			}
			this.Detach(key)
		}
	} else {
		for key, _ := range this.objs {
			if this.isInternalKey(key){
				if key==fieldWorker{
					continue
				}
				this.Detach(key)
			}
		}
	}
}

func (this *c) Print() {
	var num = 1
	fmt.Println()
	fmt.Printf("Existing components : \n")
	for key, _ := range this.objs {
		fmt.Printf("\t%d. %s \n", num, key)
		num++
	}
	fmt.Println()
}
