package mangodex

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Errors to check with errors.Is, the error types below match them when it applies.
var (
	ErrNotFound    = errors.New("not found")
	ErrRateLimited = errors.New("rate limited")
	ErrMaintenance = errors.New("under maintenance")
)

// HTTPError: The API answered with a non-200 status code (other than 429, see RateLimitError).
type HTTPError struct {
	Method     string
	Path       string
	StatusCode int
	// Message: Details sent by the API, if any.
	Message string
}

func (e *HTTPError) Error() string {
	msg := fmt.Sprintf("mangodex: %s %s: %d %s", e.Method, e.Path, e.StatusCode, strings.ToLower(http.StatusText(e.StatusCode)))
	if e.Message != "" {
		msg += ": " + e.Message
	}
	return msg
}

func (e *HTTPError) Is(target error) bool {
	switch e.StatusCode {
	case http.StatusNotFound:
		return target == ErrNotFound
	case http.StatusServiceUnavailable:
		return target == ErrMaintenance
	}
	return false
}

// RateLimitError: The API answered 429, keep sending requests and the IP gets temporarily banned.
type RateLimitError struct {
	Method     string
	Path       string
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("mangodex: %s %s: 429 too many requests (retry after %s)", e.Method, e.Path, e.RetryAfter)
}

func (e *RateLimitError) Is(target error) bool {
	return target == ErrRateLimited
}

// RetryAfterDuration: How long to wait before retrying, for callers that don't import mangodex.
func (e *RateLimitError) RetryAfterDuration() time.Duration {
	return e.RetryAfter
}
