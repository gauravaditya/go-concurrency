# Periodic worker (correct pattern)

## Requirements

- Uses time.Ticker
- Stops on context.Done()
- Ensures ticker is stopped
- No goroutine leaks
- No busy loops