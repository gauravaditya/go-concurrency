Then let's keep the momentum and treat this as the final extension of the worker-pool family before we move on.

### Your Task

Implement a bounded worker pool:

```go
type Pool struct {
    jobs chan Job
}
```

Requirements:

```go
func NewPool(workers int, queueSize int) *Pool
```

* Creates an internal buffered queue of size `queueSize`.
* Starts `workers` goroutines.

```go
func (p *Pool) Submit(job Job)
```

* Enqueues a job.
* Blocks when the queue is full.

For now:

* Ignore cancellation.
* Ignore shutdown.
* Ignore result aggregation.
* Ignore retries.

Workers can simply print:

```go
fmt.Printf("worker processed %v\n", job)
```

or sleep for a bit to simulate work.

---

### Small Experiment

Once you've implemented it, try something like:

```go
pool := NewPool(2, 5)

for i := 1; i <= 20; i++ {
    fmt.Println("submitting", i)
    pool.Submit(Job{ID: i})
    fmt.Println("submitted", i)
}
```

and make each worker take:

```go
time.Sleep(time.Second)
```

per job.

Don't analyze it yet.

Just observe:

* How many jobs get submitted immediately?
* When does the producer start blocking?
* What pattern do you see in the output?

The behavior itself is the lesson.

Post your implementation and observations, and we'll review them before deciding whether to continue deeper into queueing/backpressure or move to the next concurrency pattern.
