package mangodex

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestHTTPClientOption(t *testing.T) {
	c, err := NewDexClient(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if c.client.Timeout == 0 {
		t.Error("the default client should have a timeout")
	}

	given := &http.Client{}
	options := DefaultOptions()
	options.HTTPClient = given
	if c, err = NewDexClient(options); err != nil {
		t.Fatal(err)
	}
	if c.client != given {
		t.Error("the given client should be used")
	}
}

func TestAPIRequestInterval(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":"ok","data":[]}`))
	}))
	defer server.Close()
	c, err := NewDexClient(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	c.baseAPI = server.URL

	start := time.Now()
	for range 3 {
		if _, err := c.Manga.List(url.Values{}); err != nil {
			t.Fatal(err)
		}
	}
	if elapsed := time.Since(start); elapsed < 2*apiInterval {
		t.Errorf("3 requests took %s, wanted at least %s", elapsed, 2*apiInterval)
	}
}

// a cancelled request returns right away, not after the requests waiting before it
func TestAPIRequestIntervalCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":"ok","data":[]}`))
	}))
	defer server.Close()
	c, err := NewDexClient(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	c.baseAPI = server.URL

	if _, err := c.Manga.List(url.Values{}); err != nil {
		t.Fatal(err)
	}
	// these wait their turn, ~200ms each
	for range 3 {
		go c.Manga.List(url.Values{})
	}
	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(20*time.Millisecond, cancel)
	start := time.Now()
	if _, err := c.Manga.ListContext(ctx, url.Values{}); !errors.Is(err, context.Canceled) {
		t.Errorf("wanted context.Canceled, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 150*time.Millisecond {
		t.Errorf("the cancelled request took %s", elapsed)
	}
}

func TestMangaListWithTotal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":"ok","data":[{"id":"a","type":"manga"}],"limit":1,"offset":0,"total":42}`))
	}))
	defer server.Close()
	c, err := NewDexClient(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	c.baseAPI = server.URL

	list, total, err := c.Manga.ListWithTotal(url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "a" || total != 42 {
		t.Errorf("got %d mangas, total %d", len(list), total)
	}
}

// MangaDex allows 40 /at-home/server requests a minute, on top of the api limit
func TestAtHomeServerInterval(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":"ok","baseUrl":"https://node","chapter":{"hash":"h","data":["1.png"]}}`))
	}))
	defer server.Close()
	c, err := NewDexClient(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if c.atHomeLimiter.Limit() != rate.Every(time.Minute/40) || c.atHomeLimiter.Burst() != 1 {
		t.Errorf("got %v/s burst %d, want 40/min", c.atHomeLimiter.Limit(), c.atHomeLimiter.Burst())
	}
	c.baseAPI = server.URL
	c.atHomeLimiter = rate.NewLimiter(rate.Every(300*time.Millisecond), 1)

	start := time.Now()
	for range 3 {
		if _, err := c.AtHome.Get("id", url.Values{}); err != nil {
			t.Fatal(err)
		}
	}
	if elapsed := time.Since(start); elapsed < 600*time.Millisecond {
		t.Errorf("3 requests took %s, wanted at least 600ms", elapsed)
	}
}
