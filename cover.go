package mangodex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	CoverListPath = "cover"
)

// CoverService : Provides Cover services provided by the API.
type CoverService service

type Cover struct {
	ID            string           `json:"id"`
	Type          RelationshipType `json:"type"`
	Attributes    CoverAttributes  `json:"attributes"`
	Relationships []*Relationship  `json:"relationships"`
}

type CoverAttributes struct {
	Volume      *string `json:"volume"`
	FileName    string  `json:"fileName"`
	Description *string `json:"description"`
	Version     int     `json:"version"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Locale      string  `json:"locale"`
}

// List : Get manga covers by ID list.
// https://api.mangadex.org/docs.html#operation/get-cover
// TODO: make it generic? just accept a url.Values params?
func (s *CoverService) List(ids []string, isMangaCover bool) ([]*Cover, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = CoverListPath

	params := url.Values{}
	for _, id := range ids {
		if isMangaCover {
			params.Add("manga[]", id)
		} else {
			params.Add("ids[]", id)
		}
	}

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

	return coverArtList, err
}
