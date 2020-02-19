package C

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/zander-84/go-components/libs/ali"
	"github.com/zander-84/go-components/libs/cache"
	"github.com/zander-84/go-components/libs/cache/memory"
	"github.com/zander-84/go-components/libs/cron"
	"github.com/zander-84/go-components/libs/data/message"
	"github.com/zander-84/go-components/libs/data/queue"
	"github.com/zander-84/go-components/libs/helper"
	"github.com/zander-84/go-components/libs/logger"
	"github.com/zander-84/go-components/libs/logger/zap-log"
	"github.com/zander-84/go-components/libs/middlewares"
	"github.com/zander-84/go-components/libs/mysql"
	"github.com/zander-84/go-components/libs/mysql/grom"
	"github.com/zander-84/go-components/libs/nsq"
	"github.com/zander-84/go-components/libs/redis"
	"github.com/zander-84/go-components/libs/validate"
	"github.com/zander-84/go-components/libs/worker"
	"sync"
)

var (
	componentOnce sync.Once
	componentObj  Components
	ErrKeyExist   = errors.New("Component already exists")
)

var _ Components = componentObj

type components struct {
	c
	conf *Conf // 配置文件
}

func NewComponents(confPath string, configs ...interface{}) Components {
	componentOnce.Do(func() {
		this := new(components)
		this.objs = make(map[string]interface{}, 0)
		this.container = NewContainer()
		this.conf = this.container.Build(
			BuildConf,
			SetConfPath(confPath),
			SetData(configs...),
		).(*Conf)

		componentObj = this
		componentObj.Log()
	})

	return componentObj
}

func Get() Components {
	return componentObj
}

// 获取配置文件
func (this *components) Conf() *Conf {
	return this.conf
}

func (this *components) Helper() *CHelper.Helper {
	if this.IsExist(fieldHelper) {
		return this.Get(fieldHelper).(*CHelper.Helper)
	} else {
		helper := this.container.buildCache(
			fieldHelper,
			CHelper.BulidHelper,
			CHelper.SetConfig(this.conf.Components.Helper),
		).(*CHelper.Helper)

		this.attach(fieldHelper, helper)
		return helper
	}
}

// 获取gorm 组件
func (this *components) Mysql() CMysql.Mysql {
	if this.IsExist(fieldMysql) {
		return this.Get(fieldMysql).(CMysql.Mysql)
	} else {

		mysql := this.container.buildCache(
			fieldMysql,
			CGrom.BuildGrom,
			CGrom.SetConfig(this.conf.Components.Mysql.Gorm),
		).(CMysql.Mysql)

		this.attach(fieldMysql, mysql)
		return mysql
	}
}

// 获取日志组件  需要注意的是mysql用了gorm
func (this *components) Log() CLogger.Logger {
	if this.IsExist(fieldLog) {
		return this.Get(fieldLog).(CLogger.Logger)
	} else {
		var opts = []func(interface{}){
			CLoggerZap.SetConfig(this.conf.Components.Log.Zap),
			CLoggerZap.SetHelper(this.Helper()),
		}

		if this.conf.Components.Log.Zap.MysqlHook.Enable {
			opts = append(opts, CLoggerZap.SetGorm(this.Mysql().Obj().(*gorm.DB)))
		}
		log := this.container.buildCache(fieldLog, CLoggerZap.BuildZapLog, opts...).(CLogger.Logger)
		this.attach(fieldLog, log)
		return log
	}
}

func (this *components) Cache() CCache.Cache {
	if this.IsExist(fieldCache) {
		return this.Get(fieldCache).(CCache.Cache)
	} else {
		var cache CCache.Cache
		if this.conf.Components.Cache == "redis" {
			cache = this.Redis().Cache()
		} else {
			cache = this.Memory()
		}
		this.attach(fieldCache, cache)
		return cache
	}
}

// 获取内存存储组件
func (this *components) Memory() CCache.Cache {
	if this.IsExist(fieldMemoryCache) {
		return this.Get(fieldMemoryCache).(CCache.Cache)
	} else {
		memoryCache := this.container.buildCache(
			fieldMemoryCache,
			CMemory.BuildMemory,
			CMemory.SetConfig(this.conf.Components.Memory),
		).(CCache.Cache)

		this.attach(fieldMemoryCache, memoryCache)
		return memoryCache
	}
}

// 获取redis组件
func (this *components) Redis() *CRedis.Redis {
	if this.IsExist(fieldRedis) {
		return this.Get(fieldRedis).(*CRedis.Redis)
	} else {
		redis := this.container.buildCache(fieldRedis,
			CRedis.BuildRedis,
			CRedis.SetConfig(this.conf.Components.Redis.Cluster[0]),
		).(*CRedis.Redis)

		this.attach(fieldRedis, redis)
		return redis
	}
}

func (this *components) Validator() *CValidate.Validate {
	if this.IsExist(fieldValidator) {
		return this.Get(fieldValidator).(*CValidate.Validate)
	} else {
		validator := this.container.buildCache(
			fieldValidator,
			CValidate.BuildValidate,
			CValidate.SetConfig(this.conf.Components.Validator),
		).(*CValidate.Validate)

		this.attach(fieldValidator, validator)
		return validator
	}
}

func (this *components) Nsq() *CNsq.Nsq {
	if this.IsExist(fieldNsq) {
		return this.Get(fieldNsq).(*CNsq.Nsq)
	} else {
		nsq := this.container.buildCache(
			fieldNsq,
			CNsq.BuildNsq,
			CNsq.SetConfig(this.conf.Components.Nsq),
		).(*CNsq.Nsq)

		this.attach(fieldNsq, nsq)
		return nsq
	}
}

func (this *components) Ali() *CAli.Ali {
	if this.IsExist(fieldAli) {
		return this.Get(fieldAli).(*CAli.Ali)
	} else {
		ali := this.container.buildCache(
			fieldAli,
			CAli.BuildAli,
			CAli.SetConfig(this.conf.Components.Ali, this.Helper()),
		).(*CAli.Ali)

		this.attach(fieldAli, ali)
		return ali
	}
}

func (this *components) Cron() CCron.Crontab {
	if this.IsExist(fieldCron) {
		return this.Get(fieldCron).(CCron.Crontab)
	} else {
		cron := this.container.buildCache(
			fieldCron,
			CCron.BuildRobfigCrontab,
			CCron.SetConfig(this.conf.Components.Cron),
		).(CCron.Crontab)

		this.attach(fieldCron, cron)
		return cron
	}
}

func (this *components) Worker() CWorker.Dispatcher {
	if this.IsExist(fieldWorker) {
		return this.Get(fieldWorker).(CWorker.Dispatcher)
	} else {
		worker := this.container.buildCache(
			fieldWorker,
			CWorker.BuildDispatcher,
			CWorker.SetConfig(this.conf.Components.Worker),
		).(CWorker.Dispatcher)

		this.attach(fieldWorker, worker)
		return worker
	}
}

func (this *components) Middleware() *CMiddlewares.MiddleWares {
	if this.IsExist(fieldMiddleware) {
		return this.Get(fieldMiddleware).(*CMiddlewares.MiddleWares)
	} else {
		middleware := this.container.buildCache(
			fieldMiddleware,
			CMiddlewares.BuildMiddleWares,
		).(*CMiddlewares.MiddleWares)
		this.attach(fieldMiddleware, middleware)
		return middleware
	}
}

func (this *components) Queue(typ int) CQueue.Queue {
	return CQueue.NewQueue(typ)
}

func (this *components) Message(typ int) CMessage.Message {
	return CMessage.NewMessage(typ)
}
