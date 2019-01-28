package loadbalancer

// Job to be sent to the Worker
type Job struct {
	Fn         func() int
	ReturnChan chan int
}
