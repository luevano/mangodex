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
	// CheckIfMangaFollowedPath = "/user/follows/manga/%s"
	// ToggleMangaFollowPath    = "/manga/%s/follow"
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
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaPath, id)
	u.RawQuery = params.Encode()

	var res MangaResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Data, &manga)
	if err != nil {
		return nil, err
	}

	return manga, nil
}

// List: Get manga list.
//
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-search-manga
func (s *MangaService) List(params url.Values) (mangaList []*Manga, err error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = MangaListPath
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Data, &mangaList)
	if err != nil {
		return nil, err
	}

	return mangaList, nil
}

// TODO: update viable methods later. Most of this is either deprecated
// or the API changed drastically (due to auth being different).
// The code is heavily outdated.
/*
// CheckIfMangaFollowed : Check if a user follows a manga.
func (s *MangaService) CheckIfMangaFollowed(id string) (bool, error) {
	return s.CheckIfMangaFollowedContext(context.Background(), id)
}

// CheckIfMangaFollowedContext : CheckIfMangaFollowed with custom context.
func (s *MangaService) CheckIfMangaFollowedContext(ctx context.Context, id string) (bool, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(CheckIfMangaFollowedPath, id)

	var r Response
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ToggleMangaFollowStatus :Toggle follow status for a manga.
func (s *MangaService) ToggleMangaFollowStatus(id string, toFollow bool) (*Response, error) {
	return s.ToggleMangaFollowStatusContext(context.Background(), id, toFollow)
}

// ToggleMangaFollowStatusContext  ToggleMangaFollowStatus with custom context.
func (s *MangaService) ToggleMangaFollowStatusContext(ctx context.Context, id string, toFollow bool) (*Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(ToggleMangaFollowPath, id)

	method := http.MethodPost // To follow
	if !toFollow {
		method = http.MethodDelete // To unfollow
	}

	var r Response
	err := s.client.RequestAndDecode(ctx, method, u.String(), nil, &r)
	return &r, err
}
*/
