So let's move directly into Phase 2.

### Task 1 — Basic Worker Pool - Using forwarder pattern

Implement from scratch:

```go
func workerPool(
    workers int,
    jobs <-chan int,
) <-chan int
```

Requirements:

* Start `workers` goroutines.
* Each worker receives a number from `jobs`.
* Worker squares the number.
* Results are sent to a returned channel.
* When all workers finish:

  * result channel must be closed.
* No goroutine leaks.

Example:

```go
jobs := make(chan int)

go func() {
    defer close(jobs)

    for i := 1; i <= 5; i++ {
        jobs <- i
    }
}()

results := workerPool(3, jobs)

for r := range results {
    fmt.Println(r)
}
```

Expected output (order doesn't matter):

```text
1
4
9
16
25
```

Don't worry about cancellation, ordering, retries, buffering, or backpressure yet.

Implement it however you think is best and then we'll review it.
