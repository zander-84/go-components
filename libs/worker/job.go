package CWorker
type Job interface {
	Run() error
}


type JobFunc func() error

func (h JobFunc) Run() error {
	return h()
}
