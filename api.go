// Package mangodex provides an API wrapper for MangaDex v5.10.0 API.
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

	common       service
	refreshToken string // Unused

	// Services for MangaDex API.
	Auth            *AuthService // Deprecated
	Manga           *MangaService
	Volume          *VolumeService
	Chapter         *ChapterService
	Cover           *CoverService
	User            *UserService // Deprecated
	AtHome          *AtHomeService
	ScanlationGroup *ScanlationGroupService
}

// service: Wrapper for DexClient.
type service struct {
	client *DexClient
}

// NewDexClient: New MangaDex client.
func NewDexClient() *DexClient {
	client := http.Client{}
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	dex := &DexClient{
		client: &client,
		header: header,
	}
	dex.common.client = dex

	// Reuse the common client for the other services
	dex.Auth = (*AuthService)(&dex.common)
	dex.Manga = (*MangaService)(&dex.common)
	dex.Volume = (*VolumeService)(&dex.common)
	dex.Chapter = (*ChapterService)(&dex.common)
	dex.Cover = (*CoverService)(&dex.common)
	dex.User = (*UserService)(&dex.common)
	dex.AtHome = (*AtHomeService)(&dex.common)
	dex.ScanlationGroup = (*ScanlationGroupService)(&dex.common)

	return dex
}

// Request: Sends a request to the MangaDex API.
func (dex *DexClient) Request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = dex.header

	resp, err := dex.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		var er ErrorResponse
		var errMsg string
		err = json.NewDecoder(resp.Body).Decode(&er)
		// Sometimes the error page is just plain HTML, so it can't be decoded into ErrorResponse
		if err != nil {
			errMsg = fmt.Sprintf("failed to decode into ErrorResponse (HTML response?), error: %s", err.Error())
		} else {
			errMsg = er.GetErrors()
		}

		return nil, fmt.Errorf("non-200 status code -> (%d) %s", resp.StatusCode, errMsg)
	}

	return resp, nil
}

// RequestAndDecode: Convenience wrapper to also decode response to DexResponse type.
func (dex *DexClient) RequestAndDecode(ctx context.Context, method, url string, body io.Reader) (*DexResponse, error) {
	resp, err := dex.Request(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res DexResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
