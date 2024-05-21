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
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GetMDHomeURLPath, id)
	u.RawQuery = params.Encode()

	var res AtHomeServerResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
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
	var finalErr error
	url := strings.Join([]string{s.BaseURL, quality, s.Chapter.Hash, filename}, "/")
	ctx := context.Background()

	// Start timing how long to get all bytes for the file.
	start := time.Now()
	resp, err := s.client.Request(ctx, http.MethodGet, url, nil)
	if err != nil {
		finalErr = fmt.Errorf("Failed to get chapter page data: %s", err.Error())
	}
	defer resp.Body.Close()

	// Even on failed requests, get all the bytes from the body
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		finalErr = fmt.Errorf("Failed to read all bytes from body: %s", finalErr.Error())
	}

	// Send report in the background.
	if report {
		go func() {
			r := &reportPayload{
				URL:      url,
				Success:  finalErr == nil,
				Bytes:    len(image),
				Duration: time.Since(start).Milliseconds(),
				Cached:   strings.HasPrefix(resp.Header.Get("X-Cache"), "HIT"),
			}
			rBytes, err := json.Marshal(r)
			if err == nil {
				s.client.Request(ctx, http.MethodPost, MDHomeReportURL, bytes.NewBuffer(rBytes))
			}
		}()
	}

	if finalErr != nil {
		return nil, finalErr
	}
	return image, nil
}

// reportPayload: Required fields for reporting page download result.
type reportPayload struct {
	URL      string
	Success  bool
	Bytes    int
	Duration int64
	Cached   bool
}
