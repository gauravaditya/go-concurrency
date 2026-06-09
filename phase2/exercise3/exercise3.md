Let's take the next step.

### Task 2 — Cancellation-Aware Worker Pool

Modify the design so that the caller can cancel all workers.

Signature:

```go
func workerPool(
    ctx context.Context,
    workers int,
    jobs <-chan int,
) <-chan int
```

Requirements:

1. If `ctx.Done()` is closed:

   * workers stop immediately
   * no goroutine leaks

2. If cancellation happens halfway through:

   * unfinished jobs may be discarded

3. Results channel must still close correctly.

Example:

```go
ctx, cancel := context.WithCancel(context.Background())

results := workerPool(ctx, 3, jobs)

time.AfterFunc(time.Second, cancel)

for r := range results {
    fmt.Println(r)
}
```

Don't worry about retries, timeouts, or preserving order.

Try implementing it from scratch. After that we'll review it and discuss any subtle leak scenarios that appear.
