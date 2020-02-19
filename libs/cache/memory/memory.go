package CMemory

import (
	"github.com/patrickmn/go-cache"
	"github.com/zander-84/go-components/libs/cache"
	"reflect"
	"sync"
	"time"
)
var _ CCache.Cache = new(Memory)

type Memory struct {
	obj  *cache.Cache
	conf Conf
	once sync.Once
}

// *Memory
func NewMemory(opts ...func(interface{})) *Memory {
	_memory := &Memory{}
	for _,opt := range(opts){
		opt(_memory)
	}

	_memory.build()
	return _memory
}

func BuildMemory(opts ...func(interface{}))interface{}  {
	return NewMemory(opts...)
}
func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Memory)
		g.conf = conf
		g.conf.SetDefault()
	}
}

func (c *Memory) build() {
	c.obj = cache.New(time.Duration(c.conf.Expiration)*time.Minute, time.Duration(c.conf.CleanupInterval)*time.Minute)
}

func (c *Memory) Obj() interface{} {
	return c.obj
}

func (c *Memory) Get(key string, ptrValue interface{}) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = CCache.ErrInvalidValue
		}
	}()

	v := reflect.ValueOf(ptrValue)
	if v.Type().Kind() != reflect.Ptr {
		return CCache.ErrInvalidValue
	}

	value, found := c.obj.Get(key)
	if !found {
		return CCache.ErrRecordNotFound
	}

	if !v.Elem().CanSet() {
		err = CCache.ErrUnaddressable
		return
	}

	if reflect.ValueOf(value).Type().Kind() == reflect.Ptr {
		v.Elem().Set(reflect.ValueOf(value).Elem())
	} else {
		v.Elem().Set(reflect.ValueOf(value))
	}

	return err
}

// -1 永不过期
//______________________________________________________________________
func (c *Memory) Set(key string, value interface{}, expires time.Duration) (err error) {
	c.obj.Set(key, value, expires)

	return err
}

func (c *Memory) Delete(key string) (err error) {
	c.obj.Delete(key)
	return err
}

func (c *Memory) GetOrSet(key string, ptrValue interface{}, f func() (value interface{}, err error), expires time.Duration) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = CCache.ErrInvalidValue
		}
	}()

	getEer := c.Get(key, ptrValue)
	if getEer != nil {
		if getEer != CCache.ErrRecordNotFound {
			err = getEer
			return
		} else {
			v, ferr := f()
			if ferr != nil {
				err = ferr
				return
			}
			err = c.Set(key, v, expires)
			if err == nil {
				err = c.Get(key, ptrValue)
			}
		}
	}

	return err
}
