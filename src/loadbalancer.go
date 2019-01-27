package loadbalancer

import "container/heap"

// Pool made of a slice of workers
type Pool []*Worker

// LoadBalancer will handle the multiplexing of the jobs to workers
type LoadBalancer struct {
	workerPool Pool
	done       chan *Worker
}

func (b *LoadBalancer) balance(jobChan chan Job) {
	for {
		select {
		case todo := <-jobChan:
			b.dispatch(todo)
		case hasFinished := <-b.done:
			b.freeWorker(hasFinished)
		}
	}
}

func (b *LoadBalancer) dispatch(job Job) {
	w := heap.Pop(&b.workerPool).(*Worker)
	w.jobs <- job
	w.pending++
	heap.Push(&b.workerPool, w)
}

func (b *LoadBalancer) freeWorker(w *Worker) {
	heap.Remove(&b.workerPool, w.idx)
	w.pending--
	heap.Push(&b.workerPool, w)
}

// CreateBalancer will return a LoadBalancer instance
func CreateBalancer(numWorker int) *LoadBalancer {
	done := make(chan *Worker, numWorker)
	b := &LoadBalancer{make(Pool, 0, numWorker), done}
	for i := 0; i < numWorker; i++ {
		w := &Worker{make(chan Job), 0, i}
		heap.Push(&b.workerPool, w)
		go w.work(done)
	}
	return b
}

// Less implements heap interface
func (p Pool) Less(i, j int) bool { return p[i].pending < p[j].pending }

// Len implements heap interface
func (p Pool) Len() int { return len(p) }

func (p *Pool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].idx = i
	a[j].idx = j
}

// Push heap implementatiton
func (p *Pool) Push(x interface{}) {
	n := len(*p)
	item := x.(*Worker)
	item.idx = n
	*p = append(*p, item)
}

// Pop heap implementatiton
func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	*p = old[0 : n-1]
	return item
}
