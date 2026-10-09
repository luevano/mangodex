package mangodex

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// Tests that don't hit the API.

func TestRetryAfter(t *testing.T) {
	tests := []struct {
		header http.Header
		want   time.Duration
	}{
		{http.Header{"Retry-After": {"7"}}, 7 * time.Second},
		{http.Header{"X-Ratelimit-Retry-After": {strconv.FormatInt(time.Now().Unix()+30, 10)}}, 30 * time.Second},
		{http.Header{}, time.Second},
	}
	for _, tt := range tests {
		if got := retryAfter(tt.header); got < tt.want-time.Second || got > tt.want {
			t.Errorf("headers %v: wanted %s, got %s", tt.header, tt.want, got)
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

func TestLocalisedStringsAltTitles(t *testing.T) {
	var l LocalisedStrings
	data := `[{"en":"Solo Leveling"},{"ko":"나 혼자만 레벨업"},{"en":"Only I Level up"}]`
	if err := json.Unmarshal([]byte(data), &l); err != nil {
		t.Fatal(err)
	}
	if got := l.GetLocalString("en", false); got != "Solo Leveling" {
		t.Errorf("wanted the first english title, got %q", got)
	}
	if got := l.GetAll("en"); len(got) != 2 || got[1] != "Only I Level up" {
		t.Errorf("wanted both english titles in order, got %q", got)
	}
	if got := l.GetLocalString("fr", true); got != "Solo Leveling" {
		t.Errorf("fallback should be the first value, got %q", got)
	}
}

func TestLocalisedStringsObject(t *testing.T) {
	var l LocalisedStrings
	if err := json.Unmarshal([]byte(`{"ja-ro":"Sousou no Frieren","en":"Frieren"}`), &l); err != nil {
		t.Fatal(err)
	}
	if got := l.GetLocalString("ja-ro", false); got != "Sousou no Frieren" {
		t.Errorf("got %q", got)
	}
	// languages of an object are sorted, so the fallback is always the same
	if got := l.GetLocalString("fr", true); got != "Frieren" {
		t.Errorf("wanted a stable fallback, got %q", got)
	}
}
