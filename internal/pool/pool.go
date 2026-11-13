package pool

import (
	"sync"

	"github.com/YoumnaSalloum/glowpulse/internal/checker"
)

// job is one URL to check, with its position in the input list.
type job struct {
	index int
	url   string
}

// Run checks every URL using `workers` goroutines at the same time.
// The results come back in the same order as the input URLs.
func Run(urls []string, workers int) []checker.Result {
	if workers < 1 {
		workers = 1
	}

	jobs := make(chan job)
	results := make([]checker.Result, len(urls))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				results[j.index] = checker.Check(j.url)
			}
		}()
	}

	for i, url := range urls {
		jobs <- job{index: i, url: url}
	}
	close(jobs)

	wg.Wait()
	return results
}
