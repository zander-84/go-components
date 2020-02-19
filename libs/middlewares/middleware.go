package CMiddlewares

import (
	"github.com/zander-84/go-components/libs/middlewares/gin"
	"sync"
)

type MiddleWares struct {
	gin  *CGin.GinMiddleWare
	once sync.Once
}


func NewMiddleWares(opts ...func(interface{})) *MiddleWares {
	var _middlewares = new(MiddleWares)
	for _, opt := range opts {
		opt(_middlewares)
	}
	_middlewares.build()
	return _middlewares
}

func BuildMiddleWares(opts ...func(interface{})) interface{} {
	return NewMiddleWares(opts...)
}

func (v *MiddleWares) build()  {}

func (v *MiddleWares) Gin() *CGin.GinMiddleWare {
	if v.gin == nil {
		v.gin = CGin.New()
	}
	return v.gin
}
