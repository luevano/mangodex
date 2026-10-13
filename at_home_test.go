package mangodex

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// closeRecorder: Transport that tells when a report's response body is closed.
type closeRecorder struct{ closed chan bool }

func (c closeRecorder) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err == nil && r.Method == http.MethodPost {
		resp.Body = recordedBody{resp.Body, c.closed}
	}
	return resp, err
}

type recordedBody struct {
	io.ReadCloser
	closed chan bool
}

func (b recordedBody) Close() error {
	b.closed <- true
	return b.ReadCloser.Close()
}

// newAtHomeTestServer: MangaDex@Home node answering /data/hash/ok.png with an image and anything
// else with a 404, and the reports it gets.
func newAtHomeTestServer(t *testing.T) (*AtHomeServer, chan reportPayload, chan bool) {
	t.Helper()
	reports := make(chan reportPayload, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/report":
			var p reportPayload
			json.NewDecoder(r.Body).Decode(&p)
			reports <- p
		case "/data/hash/ok.png":
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte("image"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	closed := make(chan bool, 1)
	options := DefaultOptions()
	options.HTTPClient = &http.Client{Transport: closeRecorder{closed}}
	client, err := NewDexClient(options)
	if err != nil {
		t.Fatal(err)
	}
	client.reportURL = server.URL + "/report"
	return &AtHomeServer{client: client, BaseURL: server.URL, Chapter: ChapterData{Hash: "hash"}}, reports, closed
}

func waitReport(t *testing.T, reports chan reportPayload, closed chan bool) reportPayload {
	t.Helper()
	var r reportPayload
	select {
	case r = <-reports:
	case <-time.After(5 * time.Second):
		t.Fatal("no report")
	}
	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		t.Error("the report's response body wasn't closed")
	}
	return r
}

func TestReportSuccess(t *testing.T) {
	s, reports, closed := newAtHomeTestServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	image, err := s.GetChapterPageContext(ctx, "data", "ok.png", true)
	// the report doesn't depend on the caller's ctx
	cancel()
	if err != nil || string(image) != "image" {
		t.Fatalf("got %q, %v", image, err)
	}
	r := waitReport(t, reports, closed)
	if !r.Success || r.Bytes != len("image") || !r.Cached || r.URL != s.BaseURL+"/data/hash/ok.png" {
		t.Errorf("wrong report %+v", r)
	}
}

func TestReportFailure(t *testing.T) {
	s, reports, closed := newAtHomeTestServer(t)
	if _, err := s.GetChapterPage("data", "missing.png", true); err == nil {
		t.Fatal("wanted an error")
	}
	if r := waitReport(t, reports, closed); r.Success || r.Bytes != 0 {
		t.Errorf("wanted a failed report, got %+v", r)
	}
}

func TestNoReport(t *testing.T) {
	s, reports, _ := newAtHomeTestServer(t)
	if _, err := s.GetChapterPage("data", "ok.png", false); err != nil {
		t.Fatal(err)
	}
	select {
	case r := <-reports:
		t.Errorf("unwanted report %+v", r)
	case <-time.After(100 * time.Millisecond):
	}
}
