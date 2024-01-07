package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BaseAPI = "https://api.mangadex.org"
)

type DexResponse struct {
	Result   string          `json:"result"`
	Response string          `json:"response"`
	Data     json.RawMessage `json:"data"`
	Limit    int             `json:"limit"`
	Offset   int             `json:"offset"`
	Total    int             `json:"total"`
}

// DexClient : The MangaDex client.
type DexClient struct {
	client *http.Client
	header http.Header

	common       service
	refreshToken string

	// Services for MangaDex API
	Auth            *AuthService
	Manga           *MangaService
	Chapter         *ChapterService
	Cover           *CoverService
	User            *UserService
	AtHome          *AtHomeService
	ScanlationGroup *ScanlationGroupService
}

// service : Wrapper for DexClient.
type service struct {
	client *DexClient
}

// NewDexClient : New anonymous client. To login as an authenticated user, use DexClient.Login.
func NewDexClient() *DexClient {
	// Create client
	client := http.Client{}

	// Create header
	header := http.Header{}
	header.Set("Content-Type", "application/json") // Set default content type.

	// Create the new client
	dex := &DexClient{
		client: &client,
		header: header,
	}
	// Set the common client
	dex.common.client = dex

	// Reuse the common client for the other services
	dex.Auth = (*AuthService)(&dex.common)
	dex.Manga = (*MangaService)(&dex.common)
	dex.Chapter = (*ChapterService)(&dex.common)
	dex.Cover = (*CoverService)(&dex.common)
	dex.User = (*UserService)(&dex.common)
	dex.AtHome = (*AtHomeService)(&dex.common)
	dex.ScanlationGroup = (*ScanlationGroupService)(&dex.common)

	return dex
}

// Request : Sends a request to the MangaDex API.
func (dex *DexClient) Request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Set header for request.
	req.Header = dex.header

	// Send request.
	resp, err := dex.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		// Decode to an ErrorResponse struct.
		var er ErrorResponse
		if err = json.NewDecoder(resp.Body).Decode(&er); err != nil {
			return nil, err
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		return nil, fmt.Errorf("non-200 status code -> (%d) %s", resp.StatusCode, er.GetErrors())
	}
	return resp, nil
}

// RequestAndDecode : Convenience wrapper to also decode response to required data type
func (dex *DexClient) RequestAndDecode(ctx context.Context, method, url string, body io.Reader) (*DexResponse, error) {
	// Get the response of the request.
	resp, err := dex.Request(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Decode the request into the specified ResponseType.
	var res DexResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return &res, err
}
