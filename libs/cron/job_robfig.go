package CCron

import (
	"github.com/robfig/cron"
	"sync"
	"time"
)

type robfigCrontab struct {
	obj *cron.Cron
	conf  Conf
	jobs  map[string]*Job
	mutex sync.Mutex

}

func NewRobfigCrontab(opts ...func(interface{})) Crontab {
	var this = new(robfigCrontab)
	for _, opt := range opts {
		opt(this)
	}
	this.build()
	return this
}
func BuildRobfigCrontab(opts ...func(interface{}))interface{}  {
	return NewRobfigCrontab(opts...)
}
func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*robfigCrontab)
		g.conf = conf
		g.conf.SetDefault()
	}
}

func (this *robfigCrontab) build() {
	this.jobs = make(map[string]*Job)
	location, _ := time.LoadLocation(this.conf.Location, )
	this.obj = cron.New(cron.WithLocation(location), cron.WithSeconds())
}

func (this *robfigCrontab) AddJob(job *Job) error{
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _,ok:=this.jobs[job.ID];ok{
		return ErrIDExist
	}
	// 添加job
	id, err := this.obj.AddJob(job.Spec, job.Cmd)
	if err==nil{
		job.Obj = this.obj.Entry(id)
		this.jobs[job.ID] = job
	}
	return err
}

// 移除
func (this *robfigCrontab) Remove(id string){
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if job,ok:=this.jobs["id"];ok{
		entry:=job.Obj.(cron.Entry)
		this.obj.Remove(entry.ID)
		delete(this.jobs, id)
	}
}

//
func (this *robfigCrontab) Status() map[string]*Job{
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, job := range this.jobs {
		entry:=job.Obj.(cron.Entry)
		job.Obj = this.obj.Entry(entry.ID)
	}
	return this.jobs
}
func (this *robfigCrontab) Start(){this.obj.Start()}

func (this *robfigCrontab) Stop(){this.obj.Stop()}

func (this *robfigCrontab) Restart(id string) error{
	if job,ok:=this.jobs["id"];ok{
		this.Remove(id)
		return this.AddJob(job)
	}
	return ErrID
}

func (this *robfigCrontab) Obj() interface{}{
	return this.obj
}
