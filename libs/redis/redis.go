package CRedis

import (
	"github.com/go-redis/cache/v7"
	cache2 "github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"
	CCache "github.com/zander-84/go-components/libs/cache"
	"reflect"
	"time"
)

type Redis struct {
	obj   *redis.Client
	conf  Conf
	cache *RedisCache
}

var _ CCache.Cache = new(RedisCache)

type RedisCache struct {
	obj *cache.Codec
}

// *Redis
func NewRedis(opts ...func(interface{})) *Redis {
	var _rdb = new(Redis)
	for _, opt := range opts {
		opt(_rdb)
	}
	_rdb.build()

	return _rdb
}
func BuildRedis(opts ...func(interface{})) interface{} {
	return NewRedis(opts...)
}
func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Redis)
		g.conf = conf
		g.conf.SetDefault()
	}
}

func (c *Redis) construct(conf Conf) *Redis {
	c.conf = conf
	c.conf.SetDefault()
	c.build()
	return c
}

func (c *Redis) build() {
	c.obj = redis.NewClient(&redis.Options{
		Addr:         c.conf.Addr,
		Password:     c.conf.Password,
		DB:           c.conf.Db,
		PoolSize:     c.conf.PoolSize,
		IdleTimeout:  time.Duration(c.conf.IdleTimeout) * time.Second,
		MinIdleConns: c.conf.MinIdle,
	})
	c.cache = new(RedisCache)
	c.cache.obj = &cache.Codec{
		Redis: c.obj,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func (c *Redis) Cache() *RedisCache {
	return c.cache
}

/*_________________________________________________________________________________________*/
//-- redis action
//-- https://github.com/go-redis/redis
/*_________________________________________________________________________________________*/
func (c *Redis) Obj() *redis.Client {
	return c.obj
}
func (c *Redis) SetString(key string, str string, expires time.Duration) (err error) {
	return c.obj.Set(key, str, expires).Err()
}

func (c *Redis) GetString(key string, value interface{}) (string, error) {
	return c.obj.Get(key).Result()
}

func (c *Redis) GetBytes(key string, value interface{}) ([]byte, error) {
	return c.obj.Get(key).Bytes()
}

func (c *Redis) Dels(keys ...string) (err error) {
	return c.obj.Del(keys...).Err()
}

//-- cache
//-- https://github.com/go-redis/cache
/*_________________________________________________________________________________________*/
func (c *RedisCache) Obj() interface{} {
	return c.obj
}

func (c *RedisCache) Get(key string, ptrValue interface{}) (err error) {
	v := reflect.ValueOf(ptrValue)
	if v.Type().Kind() != reflect.Ptr {
		return CCache.ErrInvalidValue
	}

	return c.obj.Get(key, ptrValue)
}

func (c *RedisCache) Set(key string, value interface{}, expires time.Duration) (err error) {
	return c.obj.Set(&cache2.Item{
		Key:        key,
		Object:     value,
		Expiration: expires,
	})
}

func (c *RedisCache) Delete(key string) (err error) {
	return c.obj.Delete(key)
}

func (c *RedisCache) GetOrSet(key string, ptrValue interface{}, f func() (value interface{}, err error), expires time.Duration) (err error) {
	v := reflect.ValueOf(ptrValue)
	if v.Type().Kind() != reflect.Ptr {
		return CCache.ErrInvalidValue
	}

	return c.obj.Once(&cache2.Item{
		Key:        key,
		Object:     ptrValue,
		Func:       f,
		Expiration: expires,
	})
}

func (c *RedisCache) Exists(key string) bool {
	return c.obj.Exists(key)
}
