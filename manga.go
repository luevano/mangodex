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
func (m *Manga) GetTitle(langCode string) string {
	if title := m.Attributes.Title.GetLocalString(langCode); title != "" {
		return title
	}
	return m.Attributes.AltTitles.GetLocalString(langCode)
}

// GetDescription: Get description of the manga.
func (m *Manga) GetDescription(langCode string) string {
	return m.Attributes.Description.GetLocalString(langCode)
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
	State                  string             `json:"state"`
	Version                int                `json:"version"`
	CreatedAt              string             `json:"createdAt"`
	UpdatedAt              string             `json:"updatedAt"`
}

// Get: Get a manga by manga id.
//
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-manga-id
func (s *MangaService) Get(id string, params url.Values) (*Manga, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaPath, id)
	u.RawQuery = params.Encode()

	res, err := s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var manga Manga
	err = json.Unmarshal(res.Data, &manga)
	if err != nil {
		return nil, err
	}

	return &manga, nil
}

// List: Get manga list.
//
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-search-manga
func (s *MangaService) List(params url.Values) ([]*Manga, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = MangaListPath
	u.RawQuery = params.Encode()

	res, err := s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var mangaList []*Manga
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
