package C

import (
	"github.com/zander-84/go-components/libs/ali"
	"github.com/zander-84/go-components/libs/cache"
	"github.com/zander-84/go-components/libs/cron"
	"github.com/zander-84/go-components/libs/data/message"
	"github.com/zander-84/go-components/libs/data/queue"
	"github.com/zander-84/go-components/libs/helper"
	"github.com/zander-84/go-components/libs/logger"
	"github.com/zander-84/go-components/libs/middlewares"
	CMongo "github.com/zander-84/go-components/libs/mongo"
	"github.com/zander-84/go-components/libs/mysql"
	"github.com/zander-84/go-components/libs/nsq"
	CRedis "github.com/zander-84/go-components/libs/redis"
	"github.com/zander-84/go-components/libs/validate"
	"github.com/zander-84/go-components/libs/worker"
)

var (
	fieldLog         = "inner.log"
	fieldCache       = "inner.cache"
	fieldMemoryCache = "inner.memoryCache"
	fieldRedis       = "inner.redis"
	fieldMysql       = "inner.mysql"
	fieldMongo       = "inner.mongo"
	fieldValidator   = "inner.validator"
	fieldHelper      = "inner.helper"
	fieldAli         = "inner.ali"
	fieldMiddleware  = "inner.middleware"
	fieldCron        = "inner.cron"
	fieldNsq         = "inner.nsq"
	fieldWorker      = "inner.worker"
)

type Components interface {
	BasicComponents

	// 获取配置文件
	Conf() *Conf

	// 函数组件
	Helper() *CHelper.Helper

	// 获取mysql
	Mysql() CMysql.Mysql

	// 获取mongo
	Mongo() CMongo.Mongo

	// 日志
	Log() CLogger.Logger

	// 获取缓存
	Cache() CCache.Cache

	Redis() *CRedis.Redis

	// 获取内存缓存
	Memory() CCache.Cache

	// 验证组件
	Validator() *CValidate.Validate

	Nsq() *CNsq.Nsq

	//ali 组件
	Ali() *CAli.Ali

	Cron() CCron.Crontab

	Worker() CWorker.Dispatcher

	// 队列
	Queue(typ int) CQueue.Queue

	// 消息处理
	Message(typ int) CMessage.Message

	Middleware() *CMiddlewares.MiddleWares
}
