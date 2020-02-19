package CGinRateLimit

type Conf struct {
	Max          float64
	Burst        int    //实际并发
	RemoveHeader bool
}
