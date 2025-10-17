package executor

import (
	"context"
	"runtime"
	"sync"
)

// Job represents a unit of work to be processed
type Job func(ctx context.Context) error

// Pool manages a worker pool for concurrent job execution
type Pool struct {
	workers int
	jobs    chan Job
	results chan error
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewPool creates a new worker pool
func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Pool{
		workers: workers,
		jobs:    make(chan Job),
		results: make(chan error),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start begins processing jobs
func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// Submit adds a job to the pool
func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

// Stop gracefully shuts down the pool
func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}

// Results returns the results channel
func (p *Pool) Results() <-chan error {
	return p.results
}

// Cancel cancels all pending jobs
func (p *Pool) Cancel() {
	p.cancel()
}

// worker processes jobs from the queue
func (p *Pool) worker() {
	defer p.wg.Done()

	for job := range p.jobs {
		select {
		case <-p.ctx.Done():
			p.results <- p.ctx.Err()
			return
		default:
			err := job(p.ctx)
			p.results <- err
		}
	}
}
