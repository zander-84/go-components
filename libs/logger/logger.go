package CLogger

const (
	FieldRaw      = "raw"
	FieldFrom     = "from"
	FieldUid      = "uid"
	FieldDuration = "duration"
	FieldTrace    = "traceid"
	FieldSpan     = "spanid"
)

// 日志数据格式
type Data struct {
	TraceId  string      // trace id
	SpanId   string      // span id
	Uid      int         // 用户id
	Msg      string      // 消息
	Raw      interface{} // 源数据
	From     string      //来源
	Duration float64     //持续时间
}


// 日志 方法自行实现 hook
type Logger interface {

	// 简版
	Debug(data string)
	Info(data string)
	Error(data string)
	Panic(data string)
	Fatal(data string)

	// 扩展版
	ExtendDebug(data Data)
	ExtendInfo(data Data)
	ExtendError(data Data)
	ExtendPanic(data Data)
	ExtendFatal(data Data)

	// 接口对象
	Obj() interface{}
}
