package CLoggerZapMysql

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	CHelper "github.com/zander-84/go-components/libs/helper"
	"github.com/zander-84/go-components/libs/logger"
	"sync"
	"time"
)

type Log struct {
	Id       int       `json:"id,omitempty" gorm:" primary_key ;type:bigint(9) unsigned AUTO_INCREMENT; not null; comment:'主键ID' "`
	TraceId  string    `json:"trace_id,omitempty" gorm:""`
	SpanId   string    `json:"span_id,omitempty" gorm:""`
	Level    string    `json:"level,omitempty" gorm:""`
	Msg      string    `json:"msg,omitempty" gorm:""`
	From     string    `json:"from,omitempty" gorm:""`
	Duration float64   `json:"duration,omitempty" gorm:""`
	Uid      int       `json:"uid,omitempty" gorm:""`
	Raw      string    `json:"raw,omitempty" gorm:""`
	Ts       time.Time `json:"ts,omitempty" gorm:""`
}

func (Log) TableName() string {
	return "log"
}

func (this *Log) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"id":       this.Id,
		"trace_id": this.TraceId,
		"span_id":  this.SpanId,
		"level":    this.Level,
		"msg":      this.Msg,
		"from":     this.From,
		"duration": this.Duration,
		"uid":      this.Uid,
		"raw":      this.Raw,
		"ts":       this.Ts.Format("2006-01-02 15:04:05"),
	}

	return json.Marshal(data)
}

type MysqlHook struct {
	TableName string
	Gdb       *gorm.DB
	mu        sync.Mutex
	Helper    *CHelper.Helper
}

type Fields struct {
	CLogger.Data
	Level string
	Ts    string
}

/*
CREATE TABLE `log` (
  `id` bigint(9) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `trace_id` varchar(100) NOT NULL DEFAULT '' COMMENT '追踪id',
  `span_id` varchar(100) NOT NULL DEFAULT '' COMMENT '子id',
  `level` varchar(255) NOT NULL DEFAULT '' COMMENT '级别',
  `msg` varchar(255) NOT NULL DEFAULT '' COMMENT '消息',
  `from` varchar(100) NOT NULL DEFAULT '' COMMENT '来源',
  `duration` double unsigned DEFAULT NULL COMMENT '持续时间',
  `uid` bigint(9) unsigned DEFAULT NULL COMMENT '用户ID',
  `raw` text COMMENT '内容',
  `ts` timestamp NULL DEFAULT NULL COMMENT '日志记录时间',
  PRIMARY KEY (`id`),
  KEY `idx_log_span_id` (`span_id`),
  KEY `idx_log_level` (`level`),
  KEY `idx_log_from` (`from`),
  KEY `idx_log_trace_id` (`trace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 ;

*/
func (l *MysqlHook) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	db := l.Gdb

	var data Fields
	if err := json.Unmarshal(p, &data); err != nil {
		return len(p), err
	}
	ts, _ := l.Helper.TimeZone().Parse("2006-01-02 15:04:05", data.Ts)

	if data.Raw == nil {
		data.Raw = ""
	}

	var logStruct = Log{
		TraceId:  data.TraceId,
		SpanId:   data.SpanId,
		Level:    data.Level,
		Msg:      data.Msg,
		From:     data.From,
		Duration: data.Duration,
		Ts:       ts,
		Uid:      data.Uid,
		Raw:      fmt.Sprintf("%v", data.Raw),
	}

	var tableName string
	if l.TableName != "" {
		tableName = l.TableName
	} else {
		tableName = Log{}.TableName()
	}

	err = db.Table(tableName).Create(&logStruct).Error
	n = len(p)
	return n, err
}
