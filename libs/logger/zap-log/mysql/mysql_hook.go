package CLoggerZapMysql

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/zander-84/go-components/libs/logger"
	"sync"
	"time"
)

type Log struct {
	Id       int       `gorm:" primary_key ;type:bigint(9) unsigned AUTO_INCREMENT; not null; comment:'主键ID' "`
	TraceId  string    `gorm:" index;type:varchar(100); not null; default:''; comment:'追踪id';"`
	SpanId   string    `gorm:" index;type:varchar(100); not null; default:''; comment:'子id';"`
	Level    string    `gorm:" index;type:varchar(255); not null; default:''; comment:'级别';"`
	Msg      string    `gorm:" type:varchar(255); not null; default:''; comment:'消息';"`
	From     string    `gorm:" index;type:varchar(100); not null; default:''; comment:'来源';"`
	Duration float64   `gorm:" type:double unsigned; comment:'持续时间'"`
	Uid      int       `gorm:" type:bigint(9) unsigned; comment:'用户ID' "`
	Raw      string    `gorm:" type:text;  comment:'内容';"`
	Ts       time.Time `gorm:" type:timestamp; comment:'日志记录时间';"`
}

func (Log) TableName() string {
	return "log"
}

type MysqlHook struct {
	TableName string
	Gdb       *gorm.DB
	mu        sync.Mutex
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
	db.AutoMigrate(&Log{})

	var data Fields
	if err := json.Unmarshal(p, &data); err != nil {
		return len(p), err
	}
	ts, _ := time.Parse("2006-01-02 15:04:05", data.Ts)

	if data.Raw == nil{
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
