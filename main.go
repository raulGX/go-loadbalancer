package main

import (
	"math/rand"
	"runtime"
	"time"

	loadbalancer "github.com/raulgx/go-loadbalancer/src"
)

var numWorkers = runtime.GOMAXPROCS(0)

func main() {
	b := loadbalancer.CreateBalancer(numWorkers)
	jobChan := make(chan loadbalancer.Job, loadbalancer.JOB_BUFFER_LENGTH)
	returnChan := make(chan int, loadbalancer.JOB_BUFFER_LENGTH)
	for i := 0; i < loadbalancer.JOB_BUFFER_LENGTH; i++ {
		go func() {
			jobChan <- loadbalancer.Job{
				Fn: func() int {
					time.Sleep(2 * time.Millisecond)
					return rand.Intn(100)
					// return random
				}, ReturnChan: returnChan}
		}()
	}
	go b.Balance(jobChan)
	for i := 0; i < loadbalancer.JOB_BUFFER_LENGTH; i++ {
		ret := <-returnChan
		println(ret)
	}
}
