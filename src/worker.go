package loadbalancer

// Worker will receive a Job and perform it
type Worker struct {
	jobs    chan Job
	pending int
	idx     int
}

func (w *Worker) work(done chan *Worker) {
	for {
		job := <-w.jobs
		job.returnChan <- job.fn()
		done <- w
	}
}
