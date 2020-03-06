package CLoggerZap

import (
	"github.com/jinzhu/gorm"
	"github.com/zander-84/go-components/libs/helper"
	"github.com/zander-84/go-components/libs/logger"
	"github.com/zander-84/go-components/libs/logger/zap-log/email"
	"github.com/zander-84/go-components/libs/logger/zap-log/mysql"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"time"
)

var _ CLogger.Logger = new(ZapLog)

type ZapLog struct {
	obj    *zap.Logger
	conf   Conf
	gdb    *gorm.DB
	helper *CHelper.Helper
}

// comLogger.Logger
func NewZapLog(opts ...func(interface{})) CLogger.Logger {
	var _zapLog = new(ZapLog)
	for _, opt := range opts {
		opt(_zapLog)
	}
	_zapLog.build()
	return _zapLog
}

func BuildZapLog(opts ...func(interface{})) interface{} {
	return NewZapLog(opts...)
}

func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*ZapLog)
		g.conf = conf
		g.conf.SetDefault()
	}
}

func SetHelper(helper *CHelper.Helper) func(interface{}) {
	return func(i interface{}) {
		g := i.(*ZapLog)
		g.helper = helper
	}
}

func SetGorm(gbd *gorm.DB) func(interface{}) {
	return func(i interface{}) {
		g := i.(*ZapLog)
		g.gdb = gbd
	}
}

func (l *ZapLog) build() {
	logconf := l.conf
	newCore := make([]zapcore.Core, 0)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(i time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(i.Format("2006-01-02 15:04:05"))
	}
	allPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		//return lvl >= zapcore.ErrorLevel  highPriority
		//return lvl >= globalLevel && lvl < zapcore.ErrorLevel lowPriority
		return true
	})

	//____ 控制台输出
	if logconf.ConsoleHook.Enable {
		console := zapcore.Lock(os.Stdout)
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		newCore = append(newCore,
			zapcore.NewCore(consoleEncoder, console, allPriority),
		)
	}

	//____ 文件写入
	if logconf.FileHook.Enable {
		if err := l.helper.File().OpenLogFile(logconf.FileHook.Path, logconf.Name); err != nil {
			log.Fatalln("Create log error: ", err.Error())
		}

		fileHook := lumberjack.Logger{
			Filename:   l.helper.File().FullLogPath(logconf.FileHook.Path, logconf.Name), // 日志文件路径
			MaxSize:    logconf.FileHook.MaxSize,                                         // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: logconf.FileHook.MaxBackups,                                      // 日志文件最多保存多少个备份
			MaxAge:     logconf.FileHook.MaxAge,                                          // 文件最多保存多少天
			Compress:   false,                                                            // 是否压缩
		}
		fileWriter := zapcore.AddSync(&fileHook)
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		newCore = append(newCore,
			zapcore.NewCore(jsonEncoder, fileWriter, allPriority),
		)
	}

	//____ 邮件配置
	if logconf.EmailHook.Enable {
		emailhook := CLoggerZapEmail.EmailHook{
			Host:     logconf.EmailHook.Host,
			Port:     logconf.EmailHook.Port,
			User:     logconf.EmailHook.User,
			Password: logconf.EmailHook.Password,
			To:       logconf.EmailHook.To,
			Subject:  logconf.Name,
		}

		emailWriter := zapcore.AddSync(&emailhook)
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		newCore = append(newCore, zapcore.NewCore(jsonEncoder, emailWriter, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})))
	}

	//____ 写入数据库
	if logconf.MysqlHook.Enable {
		mysqlhook := CLoggerZapMysql.MysqlHook{}
		mysqlhook.TableName = logconf.MysqlHook.TableName
		mysqlhook.Gdb = l.gdb
		mysqlWriter := zapcore.AddSync(&mysqlhook)
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		newCore = append(newCore, zapcore.NewCore(jsonEncoder, mysqlWriter, allPriority))
	}
	if logconf.AddCaller {
		l.obj = zap.New(zapcore.NewTee(newCore...), zap.AddCaller())
	} else {
		l.obj = zap.New(zapcore.NewTee(newCore...))
	}
}

//  返回 (*zap.Logger)
func (l *ZapLog) Obj() interface{} {
	return l.obj
}

//*______________________________________________________________________*/
//__ 简版
/*______________________________________________________________________*/
func (l *ZapLog) Debug(data string) {
	l.obj.Debug(data)
}

func (l *ZapLog) Info(data string) {
	l.obj.Info(data)
}

func (l *ZapLog) Error(data string) {
	l.obj.Error(data)
}

func (l *ZapLog) Panic(data string) {
	l.obj.Panic(data)
}

func (l *ZapLog) Fatal(data string) {
	l.obj.Fatal(data)
}

//*______________________________________________________________________*/
//__ 来源版
/*______________________________________________________________________*/
func (l *ZapLog) DebugFrom(data string, from string) {
	l.obj.Debug(data, zap.String(CLogger.FieldFrom, from))
}

func (l *ZapLog) InfoFrom(data string, from string) {
	l.obj.Info(data, zap.String(CLogger.FieldFrom, from))
}

func (l *ZapLog) ErrorFrom(data string, from string) {
	l.obj.Error(data, zap.String(CLogger.FieldFrom, from))
}

func (l *ZapLog) PanicFrom(data string, from string) {
	l.obj.Panic(data, zap.String(CLogger.FieldFrom, from))
}

func (l *ZapLog) FatalFrom(data string, from string) {
	l.obj.Fatal(data, zap.String(CLogger.FieldFrom, from))
}

//*______________________________________________________________________*/
//__ 扩展板
/*______________________________________________________________________*/
func (l *ZapLog) ExtendDebug(data CLogger.Data) {
	l.obj.Debug(
		data.Msg,
		zap.String(CLogger.FieldTrace, data.TraceId),
		zap.String(CLogger.FieldSpan, data.SpanId),
		zap.String(CLogger.FieldFrom, data.From),
		zap.Int(CLogger.FieldUid, data.Uid),
		zap.Float64(CLogger.FieldDuration, data.Duration),
		zap.Reflect(CLogger.FieldRaw, data.Raw),
	)
}

func (l *ZapLog) ExtendInfo(data CLogger.Data) {
	l.obj.Info(
		data.Msg,
		zap.String(CLogger.FieldTrace, data.TraceId),
		zap.String(CLogger.FieldSpan, data.SpanId),
		zap.String(CLogger.FieldFrom, data.From),
		zap.Int(CLogger.FieldUid, data.Uid),
		zap.Float64(CLogger.FieldDuration, data.Duration),
		zap.Reflect(CLogger.FieldRaw, data.Raw),
	)
}

func (l *ZapLog) ExtendError(data CLogger.Data) {
	l.obj.Error(
		data.Msg,
		zap.String(CLogger.FieldTrace, data.TraceId),
		zap.String(CLogger.FieldSpan, data.SpanId),
		zap.String(CLogger.FieldFrom, data.From),
		zap.Int(CLogger.FieldUid, data.Uid),
		zap.Float64(CLogger.FieldDuration, data.Duration),
		zap.Reflect(CLogger.FieldRaw, data.Raw),
	)
}

func (l *ZapLog) ExtendPanic(data CLogger.Data) {
	l.obj.Panic(
		data.Msg,
		zap.String(CLogger.FieldTrace, data.TraceId),
		zap.String(CLogger.FieldSpan, data.SpanId),
		zap.String(CLogger.FieldFrom, data.From),
		zap.Int(CLogger.FieldUid, data.Uid),
		zap.Float64(CLogger.FieldDuration, data.Duration),
		zap.Reflect(CLogger.FieldRaw, data.Raw),
	)
}

func (l *ZapLog) ExtendFatal(data CLogger.Data) {
	l.obj.Fatal(
		data.Msg,
		zap.String(CLogger.FieldTrace, data.TraceId),
		zap.String(CLogger.FieldSpan, data.SpanId),
		zap.String(CLogger.FieldFrom, data.From),
		zap.Int(CLogger.FieldUid, data.Uid),
		zap.Float64(CLogger.FieldDuration, data.Duration),
		zap.Reflect(CLogger.FieldRaw, data.Raw),
	)
}
