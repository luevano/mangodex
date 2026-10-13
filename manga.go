package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	MangaPath     = "/manga/%s"
	MangaListPath = "/manga"
)

// MangaResponse: Manga response type, this differs from the common DexResponse type. Only used for some responses.
type MangaResponse struct {
	Result   string          `json:"result"`
	Response string          `json:"response"`
	Data     json.RawMessage `json:"data"`
}

// MangaService: Provides Manga services provided by the API.
type MangaService service

// Manga: Struct containing information on a manga.
type Manga struct {
	ID            string           `json:"id"`
	Type          RelationshipType `json:"type"`
	Attributes    MangaAttributes  `json:"attributes"`
	Relationships []*Relationship  `json:"relationships"`
}

// GetTitle: Get title of the manga.
//
// If the requested language code title is not found and fallback is true,
// the first available value is returned, else an empty string.
func (m *Manga) GetTitle(langCode string, fallback bool) string {
	if title := m.Attributes.Title.GetLocalString(langCode, fallback); title != "" {
		return title
	}
	return m.Attributes.AltTitles.GetLocalString(langCode, fallback)
}

// GetDescription: Get description of the manga.
//
// If the requested language code description is not found and fallback is true,
// the first available value is returned, else an empty string.
func (m *Manga) GetDescription(langCode string, fallback bool) string {
	return m.Attributes.Description.GetLocalString(langCode, fallback)
}

// MangaAttributes: Attributes for a manga.
type MangaAttributes struct {
	Title                  LocalisedStrings   `json:"title"`
	AltTitles              LocalisedStrings   `json:"altTitles"`
	Description            LocalisedStrings   `json:"description"`
	IsLocked               bool               `json:"isLocked"`
	Links                  LocalisedStrings   `json:"links"`
	OriginalLanguage       string             `json:"originalLanguage"`
	LastVolume             *string            `json:"lastVolume"`
	LastChapter            *string            `json:"lastChapter"`
	PublicationDemographic *Demographic       `json:"publicationDemographic"`
	Status                 *PublicationStatus `json:"status"`
	Year                   *int               `json:"year"`
	ContentRating          *ContentRating     `json:"contentRating"`
	Tags                   []*Tag             `json:"tags"`
	State                  MangaState         `json:"state"`
	Version                int                `json:"version"`
	CreatedAt              string             `json:"createdAt"`
	UpdatedAt              string             `json:"updatedAt"`
}

// Get: Get a manga by manga id.
//
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-manga-id
func (s *MangaService) Get(id string, params url.Values) (manga *Manga, err error) {
	return s.GetContext(context.Background(), id, params)
}

// GetContext: Same as Get, with a context.
func (s *MangaService) GetContext(ctx context.Context, id string, params url.Values) (manga *Manga, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = fmt.Sprintf(MangaPath, id)
	u.RawQuery = params.Encode()

	var res MangaResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("get manga %q: %w", id, err)
	}
	err = json.Unmarshal(res.Data, &manga)
	if err != nil {
		return nil, fmt.Errorf("get manga %q: %w", id, err)
	}

	return manga, nil
}

// List: Get manga list.
//
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-search-manga
func (s *MangaService) List(params url.Values) (mangaList []*Manga, err error) {
	return s.ListContext(context.Background(), params)
}

// ListContext: Same as List, with a context.
func (s *MangaService) ListContext(ctx context.Context, params url.Values) (mangaList []*Manga, err error) {
	mangaList, _, err = s.ListWithTotalContext(ctx, params)
	return mangaList, err
}

// ListWithTotal: Same as List, also returns the total number of mangas matching the params
// (for paging with offset and limit).
func (s *MangaService) ListWithTotal(params url.Values) (mangaList []*Manga, total int, err error) {
	return s.ListWithTotalContext(context.Background(), params)
}

// ListWithTotalContext: Same as ListWithTotal, with a context.
func (s *MangaService) ListWithTotalContext(ctx context.Context, params url.Values) (mangaList []*Manga, total int, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = MangaListPath
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, 0, fmt.Errorf("list manga: %w", err)
	}
	err = json.Unmarshal(res.Data, &mangaList)
	if err != nil {
		return nil, 0, fmt.Errorf("list manga: %w", err)
	}

	return mangaList, res.Total, nil
}
