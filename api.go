// Package mangodex provides an API wrapper for MangaDex v5.13.1 API.
package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	BaseAPI = "https://api.mangadex.org"
)

// apiInterval keeps to MangaDex's limit of ~5 api requests per second per IP.
const apiInterval = 200 * time.Millisecond

// maxBodySize caps the responses read: api json and page images.
const maxBodySize = 100 << 20

// DexResponse: Generic MangaDex API response type, most responses have this structure.
type DexResponse struct {
	Result   string          `json:"result"`
	Response string          `json:"response"`
	Data     json.RawMessage `json:"data"`
	Limit    int             `json:"limit"`
	Offset   int             `json:"offset"`
	Total    int             `json:"total"`
}

// DexClient: The MangaDex client.
type DexClient struct {
	client *http.Client
	header http.Header
	// BaseAPI and MDHomeReportURL, only changed in tests
	baseAPI   string
	reportURL string

	// limiter: one api request per apiInterval
	limiter *rate.Limiter

	common service

	// Services for MangaDex API.
	Manga           *MangaService
	Volume          *VolumeService
	Chapter         *ChapterService
	Cover           *CoverService
	User            *UserService
	AtHome          *AtHomeService
	ScanlationGroup *ScanlationGroupService
}

// service: Wrapper for DexClient.
type service struct {
	client *DexClient
}

// NewDexClient: New MangaDex client.
func NewDexClient(options Options) (*DexClient, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("mangodex: invalid options: %w", err)
	}

	client := options.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultTimeout}
	}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", options.UserAgent)

	dex := &DexClient{
		client:    client,
		header:    header,
		baseAPI:   BaseAPI,
		reportURL: MDHomeReportURL,
		limiter:   rate.NewLimiter(rate.Every(apiInterval), 1),
	}
	dex.common.client = dex

	// Reuse the common client for the other services
	dex.Manga = (*MangaService)(&dex.common)
	dex.Volume = (*VolumeService)(&dex.common)
	dex.Chapter = (*ChapterService)(&dex.common)
	dex.Cover = (*CoverService)(&dex.common)
	dex.User = (*UserService)(&dex.common)
	dex.AtHome = (*AtHomeService)(&dex.common)
	dex.ScanlationGroup = (*ScanlationGroupService)(&dex.common)

	return dex, nil
}

// Request: Sends a request to the MangaDex API.
func (c *DexClient) Request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("mangodex: %s %s: %w", method, url, err)
	}
	req.Header = c.header
	path := req.URL.Path

	// MangaDex@Home nodes and reports don't count
	if strings.HasPrefix(url, c.baseAPI) {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("mangodex: %s %s: %w", method, path, err)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mangodex: %s %s: %w", method, path, err)
	}
	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, &RateLimitError{Method: method, Path: path, RetryAfter: retryAfter(resp.Header)}
	}

	// the body is usually an ErrorResponse, but can also be an HTML error page
	var er ErrorResponse
	message := ""
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxBodySize)).Decode(&er); err == nil {
		message = er.GetErrors()
	}
	if message == "" && resp.StatusCode == http.StatusForbidden {
		message = "expired MangaDex@Home url, or temporarily IP banned"
	}
	return nil, &HTTPError{Method: method, Path: path, StatusCode: resp.StatusCode, Message: message}
}

// RequestAndDecode: Convenience wrapper to also decode response to given interface.
func (c *DexClient) RequestAndDecode(ctx context.Context, method, url string, body io.Reader, res any) error {
	resp, err := c.Request(ctx, method, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(io.LimitReader(resp.Body, maxBodySize)).Decode(&res); err != nil {
		return fmt.Errorf("mangodex: %s %s: decode response: %w", method, resp.Request.URL.Path, err)
	}
	return nil
}

// retryAfter: How long to wait after a 429, the endpoint-specific limits only send
// X-RateLimit-Retry-After (a unix timestamp of when the limit resets).
func retryAfter(header http.Header) time.Duration {
	if secs, err := strconv.Atoi(header.Get("Retry-After")); err == nil {
		return time.Duration(secs) * time.Second
	}
	if reset, err := strconv.ParseInt(header.Get("X-RateLimit-Retry-After"), 10, 64); err == nil {
		return max(time.Second, time.Until(time.Unix(reset, 0)).Round(time.Second))
	}
	// no hint, MangaDex's general limit is ~5 requests per second
	return time.Second
}
