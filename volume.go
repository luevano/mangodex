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
	Result  string     `json:"result"`
	Volumes VolumeList `json:"volumes"`
}

type VolumeList map[string]*Volume

func (v *VolumeList) UnmarshalJSON(data []byte) error {
	// Try to unmarshal directly into the type
	var volumes map[string]*Volume
	if err := json.Unmarshal(data, &volumes); err == nil {
		*v = volumes
		return nil
	}

	// Then try to unmarshal into a list of volumes (no volumes found)
	var noVolumes []any
	if err := json.Unmarshal(data, &noVolumes); err == nil {
		if len(noVolumes) != 0 {
			return fmt.Errorf("unexpected volume list; expected 0 volumes, got %d", len(noVolumes))
		}
		*v = map[string]*Volume{}
		return nil
	}

	return fmt.Errorf("unexpected volume list type: %s", string(data))
}

// VolumeService : Provides volume services provided by the API (manga/id/aggregate).
type VolumeService service

// Volume: Struct containing information on a volume.
type Volume struct {
	Volume   string            `json:"volume"`
	Count    int               `json:"count"`
	Chapters VolumeChapterList `json:"chapters"`
}

type VolumeChapterList map[string]VolumeChapter

func (v *VolumeChapterList) UnmarshalJSON(data []byte) error {
	// Try to unmarshal directly into the type
	var chapters map[string]VolumeChapter
	if err := json.Unmarshal(data, &chapters); err == nil {
		*v = chapters
		return nil
	}

	// Then try to unmarshal into a list of chapters (only one chapter is found)
	var oneChapter []VolumeChapter
	if err := json.Unmarshal(data, &oneChapter); err == nil {
		if len(oneChapter) != 1 {
			return fmt.Errorf("unexpected volume chapter list; expected 1 chapter, got %d", len(oneChapter))
		}
		*v = map[string]VolumeChapter{oneChapter[0].Chapter: oneChapter[0]}
		return nil
	}

	return fmt.Errorf("unexpected volume chapter list type: %s", string(data))
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

	return res.Volumes, nil
}
