At this point you've touched all the major building blocks:

✅ Worker pools

✅ Cancellation

✅ Fan-In

✅ Bounded queues

✅ Pipelines

✅ Parallel pipeline stages

That's actually most of the "core patterns" section of your roadmap.

---

The next natural step is one of two directions:

### Direction 1: Fan-Out

Take one input stream and distribute work across multiple independent downstream consumers.

This is a smaller topic because you've already done much of it implicitly with worker pools.

---

### Direction 2: Retries & Backoff

Introduce failures into the pipeline:

```text id="a2j27m"
generate
   ↓
square
   ↓
persist
```

where:

```text id="x18ddv"
persist
```

sometimes fails.

Now we need to decide:

* retry?
* how many times?
* how long to wait?
* what happens on cancellation?
* where do permanently failed jobs go?

Given everything you've completed so far, I'd actually switch to **Retries & Backoff** now. You've already experienced enough fan-out through worker pools and parallel stages that the dedicated fan-out section won't teach you much new immediately.

So next exercise would be:

> Build a worker pool where processing a job can randomly fail, and implement retry with exponential backoff and context cancellation.

That would be our first encounter with fault tolerance in concurrent systems.
