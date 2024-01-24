package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	MangaAggregatePath = "/manga/%s/aggregate"
)

// VolumeResponse: Volume response type, this differs from the common DexResponse type.
type VolumeResponse struct {
	Result  string          `json:"result"`
	Volumes json.RawMessage `json:"volumes"`
}

// VolumeService : Provides volume services provided by the API (manga/id/aggregate).
type VolumeService service

// Volume: Struct containing information on a volume.
type Volume struct {
	Volume   string                   `json:"volume"`
	Count    int                      `json:"count"`
	Chapters map[string]VolumeChapter `json:"chapters"`
}

// VolumeChapter: Chapter data specific to the volumes list. This is different to the actual Chapter data.
type VolumeChapter struct {
	Chapter string   `json:"chapter"`
	ID      string   `json:"id"`
	Others  []string `json:"others"`
	Count   int      `json:"count"`
}

// List: Get a list of manga volumes.
//
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-manga-aggregate
//
// TODO: integrate manga/id/aggregate to manga.go?
func (s *VolumeService) List(id string, params url.Values) (volumeList map[string]*Volume, err error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaAggregatePath, id)
	u.RawQuery = params.Encode()

	var res VolumeResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
	}

	// First, need to check what type is "volumes"
	// []interface = no volumes found, JSON is just an array
	// map[string]interface = volumes found, JSON is a map/dict
	var volumesType any
	err = json.Unmarshal(res.Volumes, &volumesType)
	if err != nil {
		return nil, err
	}

	switch volumesType.(type) {
	case []interface{}:
		// No volumes found case
		// Not sure how to handle the return, this is best so far
		return nil, nil
	case map[string]interface{}:
		err = json.Unmarshal(res.Volumes, &volumeList)
		if err != nil {
			return nil, err
		}
		return volumeList, nil
	default:
		return nil, fmt.Errorf("unexpected volumes response type")
	}
}
