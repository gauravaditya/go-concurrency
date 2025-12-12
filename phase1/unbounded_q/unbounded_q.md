# Unbounded Channel Using a Fan-Out Goroutine

This is the canonical Go-style unbounded queue.

Use two channels:

- A small buffered input channel where producers send
- A goroutine reading from input and writing to an internal slice queue
- A output channel for consumers

pseudocode:

```go
func NewUnboundedChan() (chan<- int, <-chan int) {
    in := make(chan int, 16)   // small buffer for producers
    out := make(chan int)      // downstream consumers
    queue := make([]int, 0)

    go func() {
        defer close(out)

        for {
            var (
                first int
                outCh chan int
            )

            if len(queue) > 0 {
                first = queue[0]
                outCh = out
            }

            select {
            case v, ok := <-in:
                if !ok {
                    // in is closed. Now drain queue before returning.
                    for _, item := range queue {
                        out <- item
                    }
                    return
                }
                queue = append(queue, v)

            case outCh <- first:
                queue = queue[1:]
            }
        }
    }()

    return in, out
}
```

## 🧠 Why this pattern is brilliant

- Prevents deadlocks
- Avoids huge channel buffers
- Allows producers to continue
- Lets you apply backpressure at the right place
- Enables batching, priority queues, network fan-out

## It's used in

- NATS
- Uber's internal Go systems
- Some parts of Kubernetes
- High-throughput schedulers
