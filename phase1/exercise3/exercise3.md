# 🚀 Exercise 3 — Go Concurrency Pipeline (Requirements)

## 🎯 Goal: Implement a 3-stage, cancellable, leak-free pipeline in Go:

[Stage 1] generator → [Stage 2] squarer → [Stage 3] printer

Each stage runs in its own goroutine and properly handles:

- context cancellation
- channel closure
- downstream propagation
- avoiding goroutine leaks

### 📦 Stage 1 — Generator

```go
func generator(ctx context.Context, nums []int) <-chan int
```

### Requirements:

- Creates an output channel.
- Sends values from nums into the channel.
- Stops immediately if ctx.Done() fires.
- Closes the output channel when done (only the producer closes).

### 📦 Stage 2 — Squarer

```go
func squarer(ctx context.Context, in <-chan int) <-chan int
```

#### Requirements: stage 2

- Reads integers from in.
- Computes the square of each number.
- Sends squared values into a new output channel.
- Stops immediately on context cancellation.
- Stops cleanly if in is closed.

Closes its output channel when it finishes.

### 📦 Stage 3 — Printer (Consumer)

```go
func printer(ctx context.Context, in <-chan int)
```

#### Requirements: stage 3

- Continuously reads from in.
- Prints each value.
- Stops immediately on context cancellation.
- Returns cleanly if the input channel closes.
- Does not close in (only producers close).

### 🔵 Pipeline Assembly (main)

The pipeline should be assembled like:

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    nums := generator(ctx, []int{1,2,3,4,5})
    squares := squarer(ctx, nums)

    printer(ctx, squares)

    fmt.Println("pipeline complete")
}
```

### ⚠ Rules You Must Follow

- Only a stage may close its own output channel.
- Never close an input channel.
- Every stage must stop immediately if ctx.Done() fires.
- Every stage must stop cleanly when upstream channels close.
- No goroutine may leak under any circumstance.
- Every send/receive involving channels must be inside select when cancellation is possible.

### 🧪 Final Expected Behavior

Running the pipeline prints:

```go-repl
1 → 1
2 → 4
3 → 9
4 → 16
5 → 25
pipeline complete
```

(Cancellation may stop early depending on your tests.)