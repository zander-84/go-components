package driver

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	obj  *mongo.Client
	conf Conf
}

func NewMongo(opts ...func(interface{})) *Mongo {
	mgo := new(Mongo)
	for _, opt := range opts {
		opt(mgo)
	}
	mgo.build()
	return mgo

}

func BuildMongo(opts ...func(interface{})) interface{} {
	return NewMongo(opts...)
}

func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Mongo)
		g.conf = conf
		g.conf.SetDefault()
	}
}

func (this *Mongo) build() {
	var err error
	dns := fmt.Sprintf("mongodb://%s:%s", this.conf.Host, this.conf.Port)
	mongoOptions := new(options.ClientOptions)
	mongoOptions.ApplyURI(dns)
	MaxPoolSize := this.conf.MaxPoolSize
	MinPoolSize := this.conf.MinPoolSize
	mongoOptions.MaxPoolSize = &MaxPoolSize
	mongoOptions.MinPoolSize = &MinPoolSize

	//mongoOptions.SetMaxConnIdleTime(time.Duration(this.conf.MaxConnIdleTime) * time.Second)

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	this.obj, err = mongo.Connect(context.TODO(), mongoOptions)
	if err != nil {
		fmt.Println("mongo Connect err: " + err.Error())
	}
	if err = this.obj.Ping(context.TODO(), nil); err != nil {
		fmt.Println("mongo err: " + err.Error())
	}

	this.obj.Database(this.conf.Database)
}

func (this *Mongo) Obj() interface{} {
	return this.obj
}

func (this *Mongo) DB() interface{} {
	return this.obj.Database(this.conf.Database)
}
