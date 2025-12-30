# Exercise A — Blocking Queue with sync.Cond

## Implement a thread-safe bounded queue using

- `sync.Mutex`
- `sync.Cond`
- No channels

This forces correct state-based waiting.

## Requirements

- Put(item): Blocks when queue is full
- Get() item: Blocks when queue is empty
- FIFO semantics
- No busy-waiting
- Correct under high concurrency
