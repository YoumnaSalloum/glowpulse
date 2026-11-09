package checker

import (
	"net/http"
	"time"
)

// Result represents the result of a URL check.
type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

// Check visits the given URL and returns a Result containing the status code, duration, and any error encountered.
func Check(url string) Result {
	client := &http.Client{Timeout: 10 * time.Second}
	start :=time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err !=nil{
		return Result{URL: url, StatusCode: 0, Duration: duration, Err: err}
	}
	defer resp.Body.Close()

	return Result{URL: url, StatusCode: resp.StatusCode, Duration: duration, Err: nil}

}
