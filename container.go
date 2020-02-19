package C

import "sync"

type Container interface {
	Build(fun func(opts ...func(interface{})) interface{}, opts ...func(interface{})) interface{}
	buildCache(key string, fun func(opts ...func(interface{})) interface{}, opts ...func(interface{})) interface{}
	delCache(key string)
	getCache(key string) (interface{}, bool)
}

var (
	containerObj     = new(container)
	containerOnce sync.Once
)

type container struct {
	objs  map[string]interface{}
	mutex sync.Mutex
}

func NewContainer() Container {
	containerOnce.Do(func() {
		containerObj.objs = make(map[string]interface{}, 0)
	})
	return containerObj
}

func (this *container) Build(newFun func(opts ...func(interface{})) interface{}, opts ...func(interface{})) interface{} {
	return newFun(opts...)
}

func (this *container) buildCache(key string, newFun func(opts ...func(interface{})) interface{}, opts ...func(interface{})) interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if res, ok := this.getCache(key); ok {
		return res
	}

	res := newFun(opts...)
	this.objs[key] = res

	return res
}

func (this *container) delCache(key string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.objs[key]; ok {
		delete(this.objs, key)
	}
}

func (this *container) getCache(key string) (interface{}, bool) {
	obj, ok := this.objs[key]
	return obj, ok
}
