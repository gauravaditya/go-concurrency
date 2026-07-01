And this brings us to a good stopping point for the first token bucket.

You've now built:

* Worker Pool
* Cancellation-aware Worker Pool
* Ordered Worker Pool
* Bounded Queue
* Fan-In
* Pipelines
* Parallel Pipeline Stages
* Retry & Backoff
* Token Bucket Rate Limiter

which is a surprisingly large chunk of the roadmap.

---

For the next exercise, I would not build another limiter immediately.

Instead I'd like to **combine** concepts.

A good challenge would be:

```text
Producer
   ↓
Bounded Queue
   ↓
Worker Pool
   ↓
Rate Limiter
   ↓
External Service
```

where:

* workers process jobs concurrently
* every request to the external service must acquire a token
* retries happen on failures
* cancellation propagates through the entire system

That would be the first exercise where several concurrency primitives interact with each other instead of being studied in isolation.

And that's usually where the most interesting bugs appear.
