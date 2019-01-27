package loadbalancer

// Job to be sent to the Worker
type Job struct {
	fn         func() int
	returnChan chan int
}
