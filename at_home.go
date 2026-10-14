package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	GetMDHomeURLPath = "/at-home/server/%s"
	MDHomeReportURL  = "https://api.mangadex.network/report"
)

// AtHomeService: Provides MangaDex@Home services provided by the API.
type AtHomeService service

// AtHomeServerResponse: A response for getting a server URL to get chapters.
type AtHomeServerResponse struct {
	Result  string      `json:"result"`
	BaseURL string      `json:"baseUrl"`
	Chapter ChapterData `json:"chapter"`
}

// ChapterData: Struct containing data for the chapter's pages.
type ChapterData struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

// AtHomeServer: Client for interfacing with MangaDex@Home.
type AtHomeServer struct {
	client *DexClient

	BaseURL string
	Chapter ChapterData
}

// Get: Get MangaDex@Home server for a chapter by id.
//
// https://api.mangadex.org/docs/redoc.html#tag/AtHome/operation/get-at-home-server-chapterId
func (s *AtHomeService) Get(id string, params url.Values) (atHome *AtHomeServer, err error) {
	return s.GetContext(context.Background(), id, params)
}

// GetContext: Same as Get, with a context.
func (s *AtHomeService) GetContext(ctx context.Context, id string, params url.Values) (atHome *AtHomeServer, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = fmt.Sprintf(GetMDHomeURLPath, id)
	u.RawQuery = params.Encode()

	if err := s.client.atHomeLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("get MangaDex@Home server for chapter %q: %w", id, err)
	}
	var res AtHomeServerResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("get MangaDex@Home server for chapter %q: %w", id, err)
	}

	atHome = &AtHomeServer{
		client:  s.client,
		BaseURL: res.BaseURL,
		Chapter: res.Chapter,
	}
	return atHome, nil
}

// GetChapterPage: Return page data for a chapter with the filename of that page.
func (s *AtHomeServer) GetChapterPage(quality, filename string, report bool) ([]byte, error) {
	return s.GetChapterPageContext(context.Background(), quality, filename, report)
}

// GetChapterPageContext: Same as GetChapterPage, with a context.
//
// With report, the result (failures included) is reported to MangaDex@Home in the background, as
// it asks clients to do so it can route around bad nodes.
func (s *AtHomeServer) GetChapterPageContext(ctx context.Context, quality, filename string, report bool) ([]byte, error) {
	url := strings.Join([]string{s.BaseURL, quality, s.Chapter.Hash, filename}, "/")

	start := time.Now()
	resp, err := s.client.Request(ctx, http.MethodGet, url, nil)
	if err != nil {
		// a cancelled download says nothing about the node
		if report && ctx.Err() == nil {
			s.report(ctx, reportPayload{URL: url, Duration: time.Since(start).Milliseconds()})
		}
		return nil, fmt.Errorf("get chapter page %q: %w", filename, err)
	}
	defer resp.Body.Close()

	image, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize+1))
	if err == nil && len(image) > maxBodySize {
		err = fmt.Errorf("image over %d bytes", maxBodySize)
	}
	if report && ctx.Err() == nil {
		s.report(ctx, reportPayload{
			URL:      url,
			Success:  err == nil,
			Bytes:    len(image),
			Duration: time.Since(start).Milliseconds(),
			Cached:   strings.HasPrefix(resp.Header.Get("X-Cache"), "HIT"),
		})
	}
	if err != nil {
		return nil, fmt.Errorf("get chapter page %q: read response: %w", filename, err)
	}
	return image, nil
}

// reportTimeout: How long a report can take, it runs in the background.
const reportTimeout = 30 * time.Second

// report: Sends the report in the background.
func (s *AtHomeServer) report(ctx context.Context, r reportPayload) {
	// the report outlives the download, so it must not be cancelled with ctx
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), reportTimeout)
	go func() {
		defer cancel()
		b, err := json.Marshal(r)
		if err != nil {
			return
		}
		resp, err := s.client.Request(ctx, http.MethodPost, s.client.reportURL, bytes.NewReader(b))
		if err == nil {
			resp.Body.Close()
		}
	}()
}

// reportPayload: Required fields for reporting page download result.
type reportPayload struct {
	URL      string `json:"url"`
	Success  bool   `json:"success"`
	Bytes    int    `json:"bytes"`
	Duration int64  `json:"duration"`
	Cached   bool   `json:"cached"`
}
