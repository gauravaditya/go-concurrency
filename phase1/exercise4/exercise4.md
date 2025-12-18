# Dynamic Fan-Out / Fan-In with Autoscaling Workers

## Goal

Build a system where:

- Jobs enter through an input channel.
- Workers dynamically scale up and down depending on load.
- Results from all workers are funneled into a single output channel.
- The whole system shuts down cleanly via context cancellation.
- No goroutine leaks.
- No dropped jobs.

This is essentially a mini autoscaling concurrency engine, like a simplified version of:

- Kafka consumer groups
- AWS Lambda concurrency manager
- Job schedulers
- High-load stream processors

This exercise forces you to understand:

- how load propagates
- how to detect pressure
- how to coordinate worker lifecycle
- how to avoid races on worker add/remove
- how to maintain fan-in correctness under scaling
- buffer dynamics
- how context cancellation cascades
- This is advanced and extremely relevant to real systems.

## 🎯 System Requirements

You will create a function

```go
func AutoScalePool(
    ctx context.Context,
    minWorkers int,
    maxWorkers int,
    jobs <-chan int,
) <-chan int
```

### Inputs

- ctx for cancellation
- a minimum number of workers to always keep alive
- a maximum number of workers allowed
- an incoming job stream

### Output

a single channel that yields all processed results

### Worker Behavior

Each worker:

- processes incoming jobs
- sends results into a fan-in channel
- exits if ctx is canceled
- exits if the autoscaler tells it to stop (scale-down event)
- must not leak
- must not block indefinitely

### Autoscaling Rules

You must implement autoscaling based on queue pressure:

Scale UP if:

- jobs are waiting in the input channel
- AND you have fewer than `maxWorkers`

Scale DOWN if:

- workers are idle for a configured interval (e.g. 100ms)
- AND current worker count > `minWorkers`

Constraints:

- Scale decisions happen in a dedicated goroutine
- Workers must not scale themselves
- Scaling logic must not race with worker exit

### Fan-In Rules

Your output channel must:

- deliver results from all workers
- close cleanly when all workers exit
- never be closed by workers directly
- be closed only after the autoscaler decides shutdown is complete

### Cancellation Behavior

When ctx is canceled:

- autoscaler must stop
- workers must stop
- output channel must close
- all goroutines must exit

cancellation cascades through multiple interacting goroutines.

### Architecture

```diagram
                 +-----------------------+
jobs --->--------|   autoscaler loop     |
                 |    decides worker     |
                 |   scale up/down       |
                 +-----------------------+
                          |  ^
         scale-up/down    |  | signals
                 v        |  |
           +--------------+--+-------------+
           |      dynamic worker set       |
           |   (size adjusts between M-N)  |
           +--------------+----------------+
                          |
               fan-in multiplexer
                          |
                          v
                  results channel
```
