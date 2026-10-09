package mangodex

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// Tests that don't hit the API.

func TestRetryAfterSecs(t *testing.T) {
	tests := []struct {
		header http.Header
		want   int
	}{
		{http.Header{"Retry-After": {"7"}}, 7},
		{http.Header{"X-Ratelimit-Retry-After": {strconv.FormatInt(time.Now().Unix()+30, 10)}}, 30},
		{http.Header{}, 1},
	}
	for _, tt := range tests {
		if got := retryAfterSecs(tt.header); got < tt.want-1 || got > tt.want {
			t.Errorf("headers %v: wanted %d, got %d", tt.header, tt.want, got)
		}
	}
}

func TestReportPayloadJSON(t *testing.T) {
	b, err := json.Marshal(reportPayload{URL: "u", Success: true, Bytes: 1, Duration: 2, Cached: true})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"url":"u","success":true,"bytes":1,"duration":2,"cached":true}`
	if string(b) != want {
		t.Errorf("wanted %s, got %s", want, b)
	}
}

func TestRelationshipEmptyAttributes(t *testing.T) {
	var r Relationship
	data := `{"id":"a1c7c817-4e59-43b7-9365-09675a149a6f","type":"author","attributes":[]}`
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		t.Errorf("empty attributes list should be ignored, got %s", err)
	}
}
