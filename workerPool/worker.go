package workerpool

import (
	"time"
)

type Worker struct {
	Id         int
	WorkerPool *WorkerPool
	quit       chan bool
}

func (w *Worker) Start() {
	for {
		select {
		// Execute the function in the Workerpool's JobQueue
		case job := <-w.WorkerPool.JobQueue:
			job.Execute()
		case <-w.quit:
			// If worker quit is true, exit loop
			return
		}
		time.Sleep(time.Duration(w.WorkerPool.WorkerCoolDown) * time.Second)
	}
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func NewWorker(id int, workerPool *WorkerPool) *Worker {
	return &Worker{Id: id, WorkerPool: workerPool, quit: make(chan bool)}
}
