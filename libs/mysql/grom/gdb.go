package CGrom

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zander-84/go-components/libs/mysql"
	"log"
	"time"
)

var _ CMysql.Mysql = new(Grom)

type Grom struct {
	obj  *gorm.DB
	conf Conf
}

func NewGrom(opts ...func(interface{})) *Grom {
	gdb := new(Grom)
	for _, opt := range opts {
		opt(gdb)
	}
	gdb.build()
	return gdb

}
func BuildGrom(opts ...func(interface{})) interface{} {
	return NewGrom(opts...)
}

func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Grom)
		g.conf = conf
		g.conf.SetDefault()
	}
}
func (this *Grom) build() {
	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", this.conf.User, this.conf.Pwd, this.conf.Host, this.conf.Port, this.conf.Database, this.conf.Charset)

	if this.obj, err = gorm.Open("mysql", dns); err != nil {
		log.Fatalln("mysql down")
	}

	this.obj.DB().SetMaxIdleConns(this.conf.MaxIdleconns)
	this.obj.DB().SetMaxOpenConns(this.conf.MaxOpenconns)
	this.obj.DB().SetConnMaxLifetime(time.Duration(this.conf.ConnMaxLifetime) * time.Second)
	this.obj.SingularTable(true)
	this.obj.LogMode(this.conf.Debug)

	// 设置时间
	this.obj.SetNowFuncOverride(func() time.Time {
		timezone := this.conf.TimeZone
		if timezone == "" {
			timezone = "Asia/Shanghai"
		}
		if location, err := time.LoadLocation(timezone); err != nil {
			log.Fatalln("mysql timezone err")
			return time.Now()
		} else {
			return time.Now().In(location)
		}
	})
}

func (this *Grom) Obj() interface{} {
	return this.obj
}

func (this *Grom) Transaction(f func(tx interface{}) (int, error)) (int, error) {
	db := this.obj
	tx := db.Begin()

	tag, e := f(tx)
	if e != nil {
		if err := tx.Rollback().Error; err != nil {
			return tag, err
		} else {
			return tag, e
		}
	} else {
		if err := tx.Commit().Error; err != nil {
			return tag, err
		} else {
			return tag, nil
		}
	}
}
