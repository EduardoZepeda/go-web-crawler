package workerpool

import (
	"sync"
	"webcrawler/config"
)

type WorkerPool struct {
	MaxWorkers int
	JobQueue   chan Job
	Wg         *sync.WaitGroup
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.MaxWorkers; i++ {
		worker := NewWorker(i, wp)
		go worker.Start()
	}
}

func (wp *WorkerPool) AddJob(job *Job) {
	wp.JobQueue <- *job
}

func NewWorkerPool(cfg *config.Config, wg *sync.WaitGroup) *WorkerPool {
	workerPool := &WorkerPool{}
	workerPool.MaxWorkers = cfg.MaxConnections
	workerPool.JobQueue = make(chan Job, len(cfg.Uris))
	workerPool.Wg = wg
	return workerPool
}
