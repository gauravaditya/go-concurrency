# Phase 2.1 — Worker Pool

## Target capabilities (what we will end with)

A worker pool that supports:

- Bounded concurrency
- Cancellation
- Graceful shutdown
- Ordered vs unordered output
- Retries with backoff (later in this phase)
- We start with the minimal correct core.