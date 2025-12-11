# A Bounded Worker Pool

This pattern appears in interviews, production systems, distributed systems, cloud runtimes, and load-balancers.
If you master this, you will dominate most Go concurrency questions.

Let’s go step-by-step.

Phase 1 — Exercise 2: Build a Bounded Worker Pool
Your goal:

## Create a system where

- The caller submits “jobs” (ints)
- A fixed number of workers (N) process these jobs
- Workers must stop cleanly when the context is canceled
- No goroutines leak
- The job channel must be closed by the caller
- Results must be collected in a result channel
- This is now a multi-component concurrency system.

### Interfaces You Must Implement

#### 1️⃣ Worker function

Same structure as Exercise 1, but now returns results:

```go
func worker(ctx context.Context, id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup)
```

#### 2️⃣ Worker Pool initializer

Write a function that:

- Spawns numWorkers workers
- Returns a jobs channel
- Returns a results channel

Example signature:

```go
func startWorkerPool(ctx context.Context, numWorkers int) (chan<- int, <-chan int)
```

#### It should

- Create a jobs channel
- Create a results channel
- Spawn workers
- Use a WaitGroup internally
- Close the results channel when all workers finish
- Return jobs and results to the caller

3️⃣ Caller code should look like this (your pool must support it):

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    jobs, results := startWorkerPool(ctx, 3)

    // send jobs
    for i := 0; i < 10; i++ {
        jobs <- i
    }
    close(jobs)

    // read results
    for r := range results {
        fmt.Println("result:", r)
    }

    fmt.Println("all workers stopped cleanly")
}
```

Expected output (order may vary):

```repl
worker 1 processed 0
worker 3 processed 1
worker 2 processed 2
...
result: 0
result: 1
...
all workers stopped cleanly
```
