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
	client *http.Client

	BaseURL string
	Chapter ChapterData
}

// Get: Get MangaDex@Home server for a chapter by id.
//
// https://api.mangadex.org/docs/redoc.html#tag/AtHome/operation/get-at-home-server-chapterId
func (s *AtHomeService) Get(id string, params url.Values) (atHome *AtHomeServer, err error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GetMDHomeURLPath, id)
	u.RawQuery = params.Encode()

	var res AtHomeServerResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
	}

	atHome = &AtHomeServer{
		client:  &http.Client{},
		BaseURL: res.BaseURL,
		Chapter: res.Chapter,
	}
	return atHome, nil
}

// GetChapterPage: Return page data for a chapter with the filename of that page.
func (s *AtHomeServer) GetChapterPage(quality, filename string) ([]byte, error) {
	path := strings.Join([]string{s.BaseURL, quality, s.Chapter.Hash, filename}, "/")
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Start timing how long to get all bytes for the file.
	start := time.Now()
	resp, err := s.client.Do(req)
	// If we cannot not get chapter page successfully.
	if err != nil || resp.StatusCode != 200 {
		if err == nil {
			err = fmt.Errorf("%d status code: failed to get chapter page data", resp.StatusCode)
		}
	}
	defer resp.Body.Close()

	var fileData []byte
	fileData, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Send report in the background.
	go func() {
		r := &reportPayload{
			URL:      path,
			Success:  err == nil,
			Bytes:    len(fileData),
			Duration: time.Since(start).Milliseconds(),
			Cached:   strings.HasPrefix(resp.Header.Get("X-Cache"), "HIT"),
		}
		_, _ = s.reportContext(ctx, r)
	}()

	return fileData, nil
}

// reportPayload: Required fields for reporting page download result.
type reportPayload struct {
	URL      string
	Success  bool
	Bytes    int
	Duration int64
	Cached   bool
}

// reportContext: Report success of getting chapter page data.
func (s *AtHomeServer) reportContext(ctx context.Context, r *reportPayload) (*http.Response, error) {
	rBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MDHomeReportURL, bytes.NewBuffer(rBytes))
	if err != nil {
		return nil, err
	}
	return s.client.Do(req)
}
