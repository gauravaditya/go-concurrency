Build a cancellable worker that does not leak goroutines.
Requirements:
✔️ 1. The worker reads jobs from a channel (jobs <-chan int)

Each job is just an integer.

✔️ 2. It must stop immediately when the context is canceled.

“Immediately" = within ≤ 50ms of calling cancel().

✔️ 3. It must NOT leak goroutines.

If:

the context is canceled, or

the job channel is closed

the worker must exit cleanly.

✔️ 4. The worker is NOT allowed to close the jobs channel.

(Only the sender owns the channel.)

✔️ 5. The worker must signal completion back to the caller via a sync.WaitGroup.