package mangodex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	CoverListPath = "/cover"
)

// CoverService: Provides Cover services provided by the API.
type CoverService service

// Cover: Struct containing information on a cover.
type Cover struct {
	ID            string           `json:"id"`
	Type          RelationshipType `json:"type"`
	Attributes    CoverAttributes  `json:"attributes"`
	Relationships []*Relationship  `json:"relationships"`
}

// CoverAttributes: Attributes for a cover.
type CoverAttributes struct {
	Volume      *string `json:"volume"`
	FileName    string  `json:"fileName"`
	Description *string `json:"description"`
	Version     int     `json:"version"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Locale      string  `json:"locale"`
}

// List : Get manga cover list.
// https://api.mangadex.org/docs/redoc.html#tag/Cover/operation/get-cover
func (s *CoverService) List(params url.Values) ([]*Cover, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = CoverListPath

	// Set query parameters
	u.RawQuery = params.Encode()

	res, err := s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var coverArtList []*Cover
	err = json.Unmarshal(res.Data, &coverArtList)
	if err != nil {
		return nil, err
	}

	return coverArtList, nil
}
