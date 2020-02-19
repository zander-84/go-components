package CCache

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidValue   = errors.New("invalid value")
	ErrUnaddressable  = errors.New("unaddressable value")
)

type Cache interface {
	// ptrValue必须是指针  获取到的是副本
	//______________________________________________________________________
	Obj() interface{}

	// ptrValue必须是指针  获取到的是副本
	//______________________________________________________________________
	Get(key string, ptrValue interface{}) error

	// 保存的结构体最好能被外部访问，私有变量在Marshal UnMarshal存在问题
	//______________________________________________________________________
	Set(key string, value interface{}, expires time.Duration) error

	Delete(key string) error

	GetOrSet(key string, ptrValue interface{}, f func() (value interface{}, err error), expires time.Duration) error
}
