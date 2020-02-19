package CWorker

type worker struct {
	jobChannel chan Job
	quit       chan bool
	running    bool // 0:等待状态  1：运行状态
	idle       bool //  0：非空闲 1：空闲
}

func newWorker(jobChannel chan Job) *worker {
	return &worker{
		jobChannel: jobChannel,
		quit:       make(chan bool),
		running:    false,
		idle:       true,
	}
}

func (this *worker) start() {
	go func() {
		for {
			this.idle = false
			this.running = false
			select {
			case job := <-this.jobChannel:
				this.running = true
				if job != nil {
					func() {
						defer func() {
							if err := recover(); err != nil {
							}
						}()
						_ = job.Run()
					}()
				}

			case <-this.quit:
				return
			}
		}
	}()

}

func (this *worker) stop() {
	this.idle = true
	this.running = false
	this.quit <- true
}
