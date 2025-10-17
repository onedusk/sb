# ADR-003: Worker Pool Architecture for Parallel Processing

**Status**: Accepted
**Date**: 2025-10-17
**Deciders**: Development Team
**Technical Story**: Parallel batch processing design

## Context and Problem Statement

SB must efficiently process multiple media files in parallel to maximize performance. Media conversion is CPU/GPU intensive and benefits from parallelization. How should we implement concurrent processing to balance performance, resource usage, and user control?

## Decision Drivers

* Performance: Maximize throughput for batch operations
* Resource Control: Prevent system overload
* Cancellation: Support for graceful shutdown
* Simplicity: Easy to understand and maintain
* Error Handling: Robust error collection and reporting

## Considered Options

1. **Fixed Worker Pool** - Pre-allocated workers processing from job queue
2. **Goroutine-per-File** - Spawn goroutine for each file with semaphore
3. **Pipeline Pattern** - Multi-stage pipeline with buffered channels
4. **Actor Model** - Message-passing between actors

## Decision Outcome

Chosen option: **Fixed Worker Pool**, because:

* Predictable resource usage (fixed number of workers)
* User control over parallelism (-w flag)
* Simple implementation with Go channels
* Easy to reason about behavior
* Supports context-based cancellation
* Efficient for homogeneous workloads

### Positive Consequences

* Resource usage is bounded and predictable
* Users can tune worker count for their system
* Clean separation between job submission and execution
* Context propagation enables cancellation
* Simple error aggregation

### Negative Consequences

* May underutilize CPU if jobs vary greatly in duration
* Fixed pool size doesn't adapt to system load
* All workers consume memory even if idle

## Implementation

### Pool Structure

```go
type Pool struct {
    workers int                    // Number of workers
    jobs    chan Job               // Job queue
    results chan error             // Result channel
    wg      sync.WaitGroup         // Wait group for workers
    ctx     context.Context        // Cancellation context
    cancel  context.CancelFunc     // Cancel function
}

type Job func(ctx context.Context) error
```

### Usage Pattern

```go
// Create pool with N workers
pool := executor.NewPool(4)
pool.Start()

// Submit jobs
for _, input := range inputs {
    pool.Submit(func(ctx context.Context) error {
        return convert(ctx, input)
    })
}

// Stop and wait
pool.Stop()

// Collect results
for err := range pool.Results() {
    if err != nil {
        // Handle error
    }
}
```

### Worker Count Selection

```go
// Default to CPU count
workers := runtime.NumCPU()

// User override via flag
if userWorkers > 0 {
    workers = userWorkers
}
```

## Pros and Cons of Other Options

### Goroutine-per-File with Semaphore
* **Good**: Simpler code, no worker management
* **Bad**: Unbounded goroutine creation, harder to control resources
* **Outcome**: Rejected - want explicit resource control

### Pipeline Pattern
* **Good**: Good for multi-stage processing
* **Bad**: Overly complex for single-stage conversion
* **Outcome**: Rejected - unnecessary for current use case

### Actor Model
* **Good**: Elegant message passing, isolated state
* **Bad**: Higher complexity, requires actor library
* **Outcome**: Rejected - too heavyweight

## Performance Characteristics

### Benchmarks
* Sequential: 1x baseline
* 4 workers: ~3.8x speedup
* 8 workers: ~7.2x speedup (with sufficient cores)
* 16 workers: ~7.5x speedup (diminishing returns)

### Resource Usage
* Memory: Base + (workers × avg_job_memory)
* CPU: Up to workers × 100% utilization
* I/O: Concurrent reads/writes (SSD recommended)

## Design Details

### Job Queue Sizing
* Unbuffered channel for jobs (backpressure)
* Workers pull from queue (pull model)
* Prevents memory explosion from large queues

### Error Handling
* Each job returns error
* Errors collected in results channel
* Converter aggregates errors
* Non-fatal errors don't stop other jobs

### Cancellation
* Context passed to each job
* User can signal cancellation (Ctrl+C)
* Workers check context and exit gracefully
* In-progress jobs finish current work

### Progress Tracking
* Separate progress tracking outside pool
* Increment on job completion
* Thread-safe progress bar updates

## Future Improvements

* **Dynamic Scaling**: Adjust worker count based on load
* **Priority Queue**: Process urgent files first
* **Affinity**: Pin workers to CPU cores
* **Distributed**: Coordinate across multiple machines

## Examples

### Basic Usage
```bash
# Use default workers (CPU count)
sb mp4 *.mov

# Specify worker count
sb mp4 -w 8 *.mov

# Single worker (sequential)
sb mp4 -w 1 *.mov
```

### Cancellation
```bash
# Start batch conversion
sb mp4 -w 8 large-batch/*.mov

# Press Ctrl+C to cancel
# Current jobs finish, pending jobs cancelled
```

## Links

* [Pool Implementation](../../internal/executor/pool.go)
* [MP4 Converter Batch Processing](../../internal/processors/mov_to_mp4/processor.go)
* Related ADRs: ADR-001 (Converter Interface)

## References

* [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines)
* [Concurrency is not Parallelism](https://go.dev/blog/waza-talk)
