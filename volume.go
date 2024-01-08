package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	MangaAggregatePath = "manga/%s/aggregate"
)

// Volume : Struct containing information on a volume.
type Volume struct {
	Volume   string                   `json:"volume"`
	Count    int                      `json:"count"`
	Chapters map[string]VolumeChapter `json:"chapters"`
}

// VolumeChapter : Chapter data specific to the volumes list. This is different to the actual Chapter data.
type VolumeChapter struct {
	Chapter string   `json:"chapter"`
	ID      string   `json:"id"`
	Others  []string `json:"others"`
	Count   int      `json:"count"`
}

// VolumeService : Provides Volume services provided by the API (manga/id/aggregate).
type VolumeService service

// List : Get a list of Manga volumes.
func (s *VolumeService) List(id string, params url.Values) (map[string]*Volume, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaAggregatePath, id)

	// Set query parameters
	u.RawQuery = params.Encode()

	res, err := s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var volumeList map[string]*Volume
	err = json.Unmarshal(res.Volumes, &volumeList)
	if err != nil {
		return nil, err
	}

	return volumeList, nil
}
