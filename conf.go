package C

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/spf13/viper"
	CAli "github.com/zander-84/go-components/libs/ali"
	CMemory "github.com/zander-84/go-components/libs/cache/memory"
	CCron "github.com/zander-84/go-components/libs/cron"
	"github.com/zander-84/go-components/libs/helper"
	comLoggerZap "github.com/zander-84/go-components/libs/logger/zap-log"
	"github.com/zander-84/go-components/libs/mysql/grom"
	CNsq "github.com/zander-84/go-components/libs/nsq"
	CRedis "github.com/zander-84/go-components/libs/redis"
	CValidate "github.com/zander-84/go-components/libs/validate"
	CWorker "github.com/zander-84/go-components/libs/worker"
	"log"
	"sync"
)

var (
	confOnce sync.Once
	confObj  = &Conf{
		obj:       viper.New(),
		BasicConf: new(BasicConf),
	}
)

type BasicConf struct {
	Components struct {
		Helper CHelper.Conf

		Mysql struct {
			Gorm CGrom.Conf
		}

		Log struct {
			Zap comLoggerZap.Conf
		}

		Cache string //memory redis   memory

		Memory CMemory.Conf

		Redis struct {
			Cluster []CRedis.Conf
		}

		Validator CValidate.Conf

		Nsq CNsq.Conf

		Ali CAli.Conf

		Cron CCron.Conf

		Worker CWorker.Conf

		Global struct {
			TimeZone string
		}
	}
}
type Conf struct {
	obj     *viper.Viper
	paths   []string
	dataPtr []interface{}
	*BasicConf
}

func BuildConf(opts ...func(interface{})) interface{} {
	return NewConf(opts...)
}

func NewConf(opts ...func(interface{})) *Conf {
	confOnce.Do(func() {

		for _, opt := range opts {
			opt(confObj)
		}

		confObj.obj.SetConfigName("config")
		confObj.obj.SetConfigType("yml")

		for _, path := range confObj.paths {
			if len(path) > 0 {
				confObj.obj.AddConfigPath(path)
			}
		}
		confObj.obj.WatchConfig()
		if err := confObj.obj.ReadInConfig(); err != nil {
			log.Fatal("配置文件加载错误" + err.Error())
		}

		// 初始化组件配置文件
		if err := confObj.obj.Unmarshal(confObj.BasicConf); err != nil {
			log.Fatal("组件配置文件反序列化错误" + err.Error())
		} else {
			fmt.Printf("%# v\n", pretty.Formatter(confObj.BasicConf))
		}

		// 初始化其它配置文件
		for _, config := range confObj.dataPtr {
			if err := confObj.obj.Unmarshal(config); err != nil {
				log.Fatal("附加配置文件加载错误" + err.Error())
			}
			fmt.Printf("%# v\n", pretty.Formatter(config))
		}

		//confObj.obj.OnConfigChange(func(in fsnotify.Event) {
		//	confObj.ReloadBasicConf()
		//})
		confObj.setGlobal()
	})
	return confObj
}

// 设置配置文件路径
func SetConfPath(confPaths []string) func(interface{}) {
	return func(i interface{}) {
		this := i.(*Conf)
		this.paths = confPaths
	}
}

// 配置初始化参数
func SetData(data ...interface{}) func(interface{}) {
	return func(i interface{}) {
		this := i.(*Conf)
		this.dataPtr = data
	}
}

// 重新加载基础配置文件
//confObj.obj.OnConfigChange(func(in fsnotify.Event) {
// if in.Name == "xxxx"{
//		confObj.ReloadBasicConf()
//	}
//})
func (this *Conf) ReloadBasicConf() {
	tmp := new(BasicConf)
	if err := this.obj.Unmarshal(tmp); err != nil {
		log.Println("组件配置文件反序列化错误" + err.Error())
	} else {
		this.BasicConf = tmp
		fmt.Printf("%# v\n", pretty.Formatter(confObj.BasicConf))
	}
	this.setGlobal()
}

func (this *Conf) Obj() interface{} {
	return this.obj
}

func (this *Conf) setGlobal() {
	timezone := this.Components.Global.TimeZone
	if timezone == "" {
		timezone = "Asia/Shanghai"
	}

	this.Components.Cron.TimeZone = timezone
	this.Components.Helper.TimeZone = timezone
	this.Components.Mysql.Gorm.TimeZone = timezone
}
