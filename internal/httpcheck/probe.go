package httpcheck

import (
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Error      string
}

func Probe(url string) Result {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	start := time.Now()

	response, err := client.Get(url)

	duration := time.Since(start)

	if err != nil {
		return Result{
			URL:      url,
			Duration: duration,
			Error:    err.Error(),
		}
	}

	defer response.Body.Close()

	return Result{
		URL:        url,
		StatusCode: response.StatusCode,
		Duration:   duration,
	}
}

func Format(result Result) string {
	if result.Error != "" {
		return fmt.Sprintf(
			"unreachable (%s)",
			result.Error,
		)
	}

	return fmt.Sprintf(
		"%d (%s)",
		result.StatusCode,
		result.Duration.Round(time.Millisecond),
	)
}
