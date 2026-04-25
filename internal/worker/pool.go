package worker

import (
	"context"
	"log"
	"sync"
)

// Job represents an asynchronous task to be executed.
type Job interface {
	Process(ctx context.Context) error
}

// Dispatcher defines the interface for dispatching jobs.
type Dispatcher interface {
	Dispatch(job Job)
}

// Pool manages a pool of worker goroutines.
type Pool struct {
	numWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup
}

// NewPool creates a new worker pool.
func NewPool(numWorkers int, queueSize int) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		jobQueue:   make(chan Job, queueSize),
	}
}

// Start spins up the background workers.
func (p *Pool) Start(ctx context.Context) {
	log.Printf("Starting worker pool with %d workers\n", p.numWorkers)

	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i+1)
	}
}

// Stop gracefully shuts down the pool.
func (p *Pool) Stop() {
	close(p.jobQueue)
	p.wg.Wait()
	log.Println("Worker pool stopped")
}

// Dispatch adds a job to the queue.
func (p *Pool) Dispatch(job Job) {
	p.jobQueue <- job
}

func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	for job := range p.jobQueue {
		if err := job.Process(ctx); err != nil {
			log.Printf("[Worker %d] Error processing job: %v\n", id, err)
		}
	}
}
