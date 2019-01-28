package loadbalancer

import "fmt"

// Worker will receive a Job and perform it
type Worker struct {
	jobs    chan Job
	pending int
	idx     int
}

func (w *Worker) work(done chan *Worker) {
	for {
		job := <-w.jobs
		fmt.Printf("Worker %v doing job, having pending %v\n", w.idx, w.pending)
		job.ReturnChan <- job.Fn()
		done <- w
	}
}
