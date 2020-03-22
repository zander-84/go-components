package CRedis

import (
	"github.com/go-redis/cache/v7"
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

func (this *Redis) construct(conf Conf) *Redis {
	this.conf = conf
	this.conf.SetDefault()
	this.build()
	return this
}

func (this *Redis) build() {
	this.obj = redis.NewClient(&redis.Options{
		Addr:         this.conf.Addr,
		Password:     this.conf.Password,
		DB:           this.conf.Db,
		PoolSize:     this.conf.PoolSize,
		IdleTimeout:  time.Duration(this.conf.IdleTimeout) * time.Second,
		MinIdleConns: this.conf.MinIdle,
	})
	this.cache = new(RedisCache)
	this.cache.obj = &cache.Codec{
		Redis: this.obj,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func (this *Redis) Cache() *RedisCache {
	return this.cache
}

/*_________________________________________________________________________________________*/
//-- redis action
//-- https://github.com/go-redis/redis
/*_________________________________________________________________________________________*/
func (this *Redis) Obj() *redis.Client {
	return this.obj
}
func (this *Redis) SetString(key string, str string, expires time.Duration) (err error) {
	return this.obj.Set(key, str, expires).Err()
}

func (this *Redis) GetString(key string, value interface{}) (string, error) {
	return this.obj.Get(key).Result()
}

func (this *Redis) GetBytes(key string, value interface{}) ([]byte, error) {
	return this.obj.Get(key).Bytes()
}

func (this *Redis) Dels(keys ...string) (err error) {
	return this.obj.Del(keys...).Err()
}

func (this *Redis) TryLockWithTimeout(identify string, duration time.Duration) (bool, error) {
	return this.obj.SetNX(identify, true, duration).Result()
}

func (this *Redis) TryLockWithWaiting(identify string, duration time.Duration, waitTime int) (bool, error) {
	for i := 0; i < waitTime; i++ {
		ok, err := this.obj.SetNX(identify, true, duration).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		time.Sleep(time.Second)
	}
	return false, nil
}

//-- cache
//-- https://github.com/go-redis/cache
/*_________________________________________________________________________________________*/
func (this *RedisCache) Obj() interface{} {
	return this.obj
}

func (this *RedisCache) Get(key string, ptrValue interface{}) (err error) {
	v := reflect.ValueOf(ptrValue)
	if v.Type().Kind() != reflect.Ptr {
		return CCache.ErrInvalidValue
	}

	return this.obj.Get(key, ptrValue)
}

func (this *RedisCache) Set(key string, value interface{}, expires time.Duration) (err error) {
	return this.obj.Set(&cache.Item{
		Key:        key,
		Object:     value,
		Expiration: expires,
	})
}

func (this *RedisCache) Delete(key string) (err error) {
	return this.obj.Delete(key)
}

func (this *RedisCache) GetOrSet(key string, ptrValue interface{}, f func() (value interface{}, err error), expires time.Duration) (err error) {
	v := reflect.ValueOf(ptrValue)
	if v.Type().Kind() != reflect.Ptr {
		return CCache.ErrInvalidValue
	}

	return this.obj.Once(&cache.Item{
		Key:        key,
		Object:     ptrValue,
		Func:       f,
		Expiration: expires,
	})
}

func (this *RedisCache) Exists(key string) bool {
	return this.obj.Exists(key)
}
