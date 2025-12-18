# go-concurrency

refresh go concurrency concepts and patterns

## How to run an exercise code

```make
make run/phase1/exercise1 #executes the main.go for phase1 -> exercise1
```

## 🚀 Golang Concurrency Mastery Roadmap

This is NOT a beginner roadmap.

### Phase 1 — Master the Foundations (1–2 weeks)

✅ 1. Goroutine Lifecycle & Scheduling

- M:N scheduler (G–P–M model)
- How goroutines are parked/unparked
- Blocking vs non-blocking operations

➤ Goal: Explain scheduler behavior in interviews & reason about CPU-bound vs IO-bound workloads.

✅ 2. Channels Deep Dive

- Buffered vs unbuffered semantics
- Channel closing protocol
- Nil channels & disabling select cases
- How to avoid channel leaks

➤ Goal: Be able to reason from first principles about deadlocks.

✅ 3. sync Package Mastery

- WaitGroup (correct usage patterns)
- Mutex vs RWMutex vs sync.Map vs atomics
- sync.Once, sync.Cond, sync.Pool

➤ Goal: Choose exactly the right primitive for any scenario.

### Phase 2 — Build Core Concurrency Patterns (2–3 weeks)

Each pattern must be implemented from scratch, tested, and benchmarked.

🔹 4. Worker Pool (bounded concurrency)

Variations:

- With result aggregation
- With cancellation
- With retries & backoff
- Ordered vs unordered outputs

🔹 5. Fan-Out / Fan-In

- Load distribution
- Fair vs unfair scheduling
- How to prevent workers from leaking

🔹 6. Pipeline Architectures

- Multi-stage processing
- Backpressure handling
- Propagating cancellation across stages

🔹 7. Context-Aware Systems

- Graceful shutdown
- Timeouts/deadlines
- Cleaning up child goroutines

🔹 8. Timers & Tickers (correct usage)

- Common ticker leaks
- Coordinating periodic workers safely

### Phase 3 — Advanced Real-World Systems (3–4 weeks)

🔥 9. Concurrency for IO-Bound Systems

- Producer–consumer queues
- Stream processing
- Designing for throughput vs latency

🔥 10. Concurrency for CPU-Bound Systems

- Sharded locking
- Work stealing
- Task partitioning
- Parallel map-reduce

🔥 11. Rate Limiting / Throttling / Debouncing

- Implement:
- Token bucket
- Leaky bucket
- Sliding window limiter
(All in idiomatic Go)

🔥 12. Building Reliable Systems

- Detecting & preventing goroutine leaks
- Memory visibility & happens-before relations
- Race detector — what it really detects

### Phase 4 — Expert Level (4–6 weeks)

🔮 13. Mastering the Go Runtime

- How the scheduler handles syscalls
- Work-stealing queues
- Preemption & cooperative scheduling
- Stack growth mechanics

🔮 14. High-Performance Concurrency

- Lock-free algorithms using atomics
- False sharing & cache line alignment
- Benchmarking with go test -bench
- Profiling goroutine contention

🔮 15. Concurrency in Distributed Systems

- Designing concurrent services (Kafka-style consumers)
- Backpressure propagation patterns
- Idempotency considerations
- Dead-letter queues

### Phase 5 — Capstone Projects (ongoing)

You will build:

- A production-grade worker pool library
- A parallel map-reduce engine
- Your own rate limiter library
- A streaming pipeline with backpressure
- A miniature task scheduler (like a tiny Kubernetes)
- I will review your code, identify flaws, and help you iterate
