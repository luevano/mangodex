package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	GetMDHomeURLPath = "/at-home/server/%s"
	MDHomeReportURL  = "https://api.mangadex.network/report"
)

// AtHomeService: Provides MangaDex@Home services provided by the API.
type AtHomeService service

// MDHomeServerResponse: A response for getting a server URL to get chapters.
type MDHomeServerResponse struct {
	Result  string      `json:"result"`
	BaseURL string      `json:"baseUrl"`
	Chapter ChapterData `json:"chapter"`
}

func (r *MDHomeServerResponse) GetResult() string {
	return r.Result
}

// ChapterData: Struct containing data for the chapter's pages.
type ChapterData struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

// MDHomeClient: Client for interfacing with MangaDex@Home.
type MDHomeClient struct {
	client  *http.Client
	quality string
	BaseURL string
	Hash    string
	Pages   []string
}

// NewMDHomeClient: Get MangaDex@Home client for a chapter.
// https://api.mangadex.org/docs/redoc.html#tag/AtHome/operation/get-at-home-server-chapterId
func (s *AtHomeService) NewMDHomeClient(chapterID string, quality string, forcePort443 bool) (*MDHomeClient, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GetMDHomeURLPath, chapterID)

	// Set query parameters
	q := u.Query()
	q.Set("forcePort443", strconv.FormatBool(forcePort443))
	u.RawQuery = q.Encode()

	resp, err := s.client.Request(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var res MDHomeServerResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if err != nil {
		return nil, err
	}

	// Set the required pages data for required quality.
	pages := res.Chapter.Data
	if quality == "data-saver" {
		pages = res.Chapter.DataSaver
	}

	return &MDHomeClient{
		client:  &http.Client{},
		quality: quality,
		BaseURL: res.BaseURL,
		Hash:    res.Chapter.Hash,
		Pages:   pages,
	}, nil
}

// GetChapterPage: Return page data for a chapter with the filename of that page.
func (c *MDHomeClient) GetChapterPage(filename string) ([]byte, error) {
	path := strings.Join([]string{c.BaseURL, c.quality, c.Hash, filename}, "/")
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Start timing how long to get all bytes for the file.
	start := time.Now()
	resp, err := c.client.Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var fileData []byte
	// If we cannot not get chapter successfully
	if err != nil || resp.StatusCode != 200 {
		if err == nil {
			err = fmt.Errorf("%d status code", resp.StatusCode)
		}
	}

	// Read file data.
	fileData, err = io.ReadAll(resp.Body)

	// Send report in the background.
	go func() {
		// Create the payload to send.
		r := &reportPayload{
			URL:      path,
			Success:  err == nil,
			Bytes:    len(fileData),
			Duration: time.Since(start).Milliseconds(),
			Cached:   strings.HasPrefix(resp.Header.Get("X-Cache"), "HIT"),
		}

		_, _ = c.reportContext(ctx, r) // Send report
	}()

	return fileData, err
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
func (c *MDHomeClient) reportContext(ctx context.Context, r *reportPayload) (*http.Response, error) {
	rBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MDHomeReportURL, bytes.NewBuffer(rBytes))
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}
