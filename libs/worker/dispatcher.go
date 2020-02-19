package CWorker

import (
	"github.com/zander-84/go-components/libs/data/queue"
	"sync/atomic"
	"time"
)

type Dispatcher interface {
	AddJob(job Job, priority int)
	Stop()
	Workers() int64
}

type dispatcher struct {
	worker          []*worker    // 工人
	idleWorkers     int64        // 空闲工人
	allowIdleNum    int64        // 允许最大空闲数
	optimizeWorkers int32        // 空值增加或者减少工人  用int代替bool
	limitQ          chan bool     // 最大限制
	priorityQ       CQueue.Queue // 有限队列
	dispatchQ       chan Job     // 真正任务派发队列
	conf            Conf
}

func NewDispatcher(opts ...func(interface{})) Dispatcher {
	var this = new(dispatcher)
	for _, opt := range opts {
		opt(this)
	}
	this.build()
	return this
}

func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*dispatcher)
		g.conf = conf
		g.conf.SetDefault()
	}
}

func (this *dispatcher) build() {
	this.worker = make([]*worker, this.conf.MaxWorkers)
	this.dispatchQ = make(chan Job, this.conf.MaxWorkers)
	this.limitQ = make(chan bool, this.conf.MaxQueues)
	this.priorityQ = CQueue.NewQueue(CQueue.MaxPriorityQ)
	this.idleWorkers = 0
	this.optimizeWorkers = 0
	this.allowIdleNum = int64(this.conf.MaxWorkers - this.conf.MinWorkers)
	this.run()
}

func BuildDispatcher(opts ...func(interface{})) interface{} {
	return NewDispatcher(opts...)
}

func (this *dispatcher) run() {
	for i := 0; i < len(this.worker); i++ {
		this.worker[i] = newWorker(this.dispatchQ)
		this.worker[i].start()
	}

	go this.dispatch()
}

func (this *dispatcher) dispatch() {
	for {
		select {

		case _ = <-this.limitQ:
			this.dispatchQ <- this.priorityQ.Pop().(Job)

			// 增加工人  任务比工人要多
			if this.idleWorkers > 0 {
				workers := this.conf.MaxQueues - int(this.idleWorkers)
				lenQ := len(this.limitQ)
				if lenQ > workers && lenQ > 2*workers {
					go this.addWorker()
				}
			}

		//一分钟 后减少 工人数
		case <-time.After(time.Minute):
			go this.reduceWorker()
		}
	}
}

func (this *dispatcher) AddJob(job Job, priority int) {
	this.priorityQ.PushPriority(job, priority)
	this.limitQ <- true
}

func (this *dispatcher) Stop() {
	for _, worker := range this.worker {
		if !worker.idle {
			worker.stop()
		}
	}
}

func (this *dispatcher) Workers() int64 {
	return int64(this.conf.MaxWorkers) - this.idleWorkers
}

// 减少工人
func (this *dispatcher) reduceWorker() {
	// 增加或者减少 同时只能有一个操作
	if atomic.AddInt32(&this.optimizeWorkers, 1) != 1 {
		return
	}
	// 还未到达最大空闲数
	if this.allowIdleNum > this.idleWorkers {

		var allowIdleNum = this.allowIdleNum
		for _, worker := range this.worker {

			// 当前空闲数量 等于 最大空闲数量 跳出
			if this.idleWorkers >= allowIdleNum {
				break
			}

			// 停止没有运行的 非空闲的线程
			if !worker.running && !worker.idle {
				// 同步关闭
				worker.stop()
				atomic.AddInt64(&this.idleWorkers, 1)
				allowIdleNum--
			}
			if allowIdleNum == 0 {
				break
			}
		}
	}
	atomic.StoreInt32(&this.optimizeWorkers, 0)
}

// 增加工人
func (this *dispatcher) addWorker() {
	// 增加或者减少 同时只能有一个操作
	if atomic.AddInt32(&this.optimizeWorkers, 1) != 1 {
		return
	}

	// 计算需要几个线程
	var length = int64(len(this.limitQ))
	var needs int64 = 0
	if this.idleWorkers > length {
		needs = length
	} else {
		needs = this.idleWorkers
	}

	// 开始寻找空闲
	for _, worker := range this.worker {
		if worker.idle {
			// 开启
			worker.start()
			atomic.AddInt64(&this.idleWorkers, -1)
			needs--
		}
		if needs == 0 {
			break
		}
	}
	atomic.StoreInt32(&this.optimizeWorkers, 0)
}
