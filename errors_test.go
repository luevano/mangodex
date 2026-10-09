package mangodex

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// newErrorTestClient: Client pointing to a local server that answers every request with the given
// status code, headers and body.
func newErrorTestClient(t *testing.T, status int, headers map[string]string, body string) *DexClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	c, err := NewDexClient(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	c.baseAPI = server.URL
	return c
}

func TestHTTPError(t *testing.T) {
	notFound := `{"result":"error","errors":[{"status":404,"title":"not_found_http_exception","detail":"Manga could not be found"}]}`
	tests := []struct {
		status   int
		body     string
		sentinel error
		message  string
	}{
		{http.StatusNotFound, notFound, ErrNotFound,
			`get manga "x": mangodex: GET /manga/x: 404 not found: not_found_http_exception: Manga could not be found`},
		{http.StatusServiceUnavailable, "<html>maintenance</html>", ErrMaintenance,
			`get manga "x": mangodex: GET /manga/x: 503 service unavailable`},
		{http.StatusForbidden, "<html>forbidden</html>", nil,
			`get manga "x": mangodex: GET /manga/x: 403 forbidden: expired MangaDex@Home url, or temporarily IP banned`},
	}
	for _, tt := range tests {
		c := newErrorTestClient(t, tt.status, nil, tt.body)

		_, err := c.Manga.Get("x", url.Values{})
		if err == nil || err.Error() != tt.message {
			t.Errorf("wanted %q, got %v", tt.message, err)
		}
		var httpErr *HTTPError
		if !errors.As(err, &httpErr) || httpErr.StatusCode != tt.status {
			t.Errorf("wanted an HTTPError with status %d, got %v", tt.status, err)
		}
		if tt.sentinel != nil && !errors.Is(err, tt.sentinel) {
			t.Errorf("%v should match %v", err, tt.sentinel)
		}
	}
}

func TestRateLimitError(t *testing.T) {
	c := newErrorTestClient(t, http.StatusTooManyRequests, map[string]string{"Retry-After": "30"}, "")

	_, err := c.Chapter.List(url.Values{})
	want := `list chapters: mangodex: GET /chapter: 429 too many requests (retry after 30s)`
	if err == nil || err.Error() != want {
		t.Fatalf("wanted %q, got %v", want, err)
	}
	if !errors.Is(err, ErrRateLimited) {
		t.Errorf("%v should match ErrRateLimited", err)
	}
	// how mangal checks it without importing mangodex
	var retry interface{ RetryAfterDuration() time.Duration }
	if !errors.As(err, &retry) || retry.RetryAfterDuration() != 30*time.Second {
		t.Errorf("wanted a 30s retry, got %v", err)
	}
}

func TestContextCanceled(t *testing.T) {
	c := newErrorTestClient(t, http.StatusOK, nil, "{}")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.Manga.GetContext(ctx, "x", url.Values{}); !errors.Is(err, context.Canceled) {
		t.Errorf("wanted context.Canceled, got %v", err)
	}
}

func TestInvalidOptions(t *testing.T) {
	if _, err := NewDexClient(Options{}); err == nil || err.Error() != "mangodex: invalid options: UserAgent is empty" {
		t.Errorf("wanted an invalid options error, got %v", err)
	}
}
