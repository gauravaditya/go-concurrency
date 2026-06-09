## Next Challenge: Ordered Worker Pool

The worker pool you've built has this property:

```text
Input:
1 2 3 4 5

Output:
1 9 4 25 16
```

Order is not guaranteed.

Many real systems require:

```text
Input:
1 2 3 4 5

Output:
1 4 9 16 25
```

while still processing jobs concurrently.

### Your task

Design (don't code immediately) a worker pool with:

* N workers
* concurrent processing
* output order identical to input order

Questions to think about:

1. How will a worker know the original position of a job?
2. If job #5 finishes before job #2, where will you store it?
3. Who is responsible for releasing results in order?
4. What happens if job #2 is very slow and jobs #3–#100 are already complete?

Describe your design first. Don't worry about perfect code yet. The design discussion is where most of the learning is.
