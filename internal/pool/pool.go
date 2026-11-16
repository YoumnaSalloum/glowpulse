package pool

import (
	"context"
	"sync"
	"time"

	"github.com/YoumnaSalloum/glowpulse/internal/checker"
)

// job is one URL to check, with its position in the input list.
type job struct {
	index int
	url   string
}

// Run checks every URL using `workers` goroutines at the same time.
// If ctx is cancelled (for example with Ctrl+C), no new checks start,
// and URLs that were never checked get a "cancelled" result.
func Run(ctx context.Context, urls []string, workers int, timeout time.Duration) []checker.Result {
	if workers < 1 {
		workers = 1
	}

	jobs := make(chan job)
	results := make([]checker.Result, len(urls))
	checked := make([]bool, len(urls))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				results[j.index] = checker.Check(ctx, j.url, timeout)
				checked[j.index] = true
			}
		}()
	}

send:
	for i, url := range urls {
		select {
		case jobs <- job{index: i, url: url}:
		case <-ctx.Done():
			break send
		}
	}
	close(jobs)
	wg.Wait()

	for i, url := range urls {
		if !checked[i] {
			results[i] = checker.Result{URL: url, Status: checker.StatusDown, Err: ctx.Err()}
		}
	}
	return results
}
