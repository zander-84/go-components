package CLoggerZapMongo

import (
	"context"
	"encoding/json"
	"fmt"
	CHelper "github.com/zander-84/go-components/libs/helper"
	"github.com/zander-84/go-components/libs/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

var Helper *CHelper.Helper

type Log struct {
	TraceId  string
	SpanId   string
	Level    string
	Msg      string
	From     string
	Tag      string
	Duration float64
	Uid      int
	Raw      string
	Ts       int64
}

//db.log.createIndex({'ts':-1})
//db.log.createIndex({'traceid':-1})
//db.log.createIndex({'spanid':-1})
//db.log.createIndex({'level':-1})
//db.log.createIndex({'from':-1})
//db.log.createIndex({'tag':-1})
//db.log.createIndex({'uid':-1})

func (Log) TableName() string {
	return "log"
}

func (this *Log) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"trace_id": this.TraceId,
		"span_id":  this.SpanId,
		"level":    this.Level,
		"msg":      this.Msg,
		"from":     this.From,
		"tag":      this.Tag,
		"duration": this.Duration,
		"uid":      this.Uid,
		"raw":      this.Raw,
		"ts":       Helper.TimeZone().LocationNow(time.Unix(this.Ts, 0)).Format("2006-01-02 15:04:05"),
	}

	return json.Marshal(data)
}

type MongoHook struct {
	TableName string
	Mdb       *mongo.Database
	mu        sync.Mutex
	Helper    *CHelper.Helper
}

type Fields struct {
	CLogger.Data
	Level string
	Ts    string
}

func (l *MongoHook) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var data Fields
	if err := json.Unmarshal(p, &data); err != nil {
		return len(p), err
	}

	if data.Raw == nil {
		data.Raw = ""
	}
	ts, _ := l.Helper.TimeZone().Parse("2006-01-02 15:04:05", data.Ts)

	var logStruct = Log{
		TraceId:  data.TraceId,
		SpanId:   data.SpanId,
		Level:    data.Level,
		Msg:      data.Msg,
		From:     data.From,
		Tag:      data.Tag,
		Duration: data.Duration,
		Ts:       ts.Unix(),
		Uid:      data.Uid,
		Raw:      fmt.Sprintf("%v", data.Raw),
	}

	var tableName string
	if l.TableName != "" {
		tableName = l.TableName
	} else {
		tableName = Log{}.TableName()
	}

	_, err = l.Mdb.Collection(tableName).InsertOne(context.TODO(), &logStruct)
	n = len(p)
	return n, err
}
