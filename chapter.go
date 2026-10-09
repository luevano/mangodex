package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	ChapterPath       = "/chapter/%s"
	ChapterListPath   = "/chapter"
	MangaChaptersPath = "/manga/%s/feed" // TODO: move to manga.go?
	// MangaReadMarkersPath = "/manga/%s/read"
)

// ChapterService: Provides chapter services provided by the API.
type ChapterService service

// Chapter: Struct containing information on a chapter.
type Chapter struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Attributes    ChapterAttributes `json:"attributes"`
	Relationships []*Relationship   `json:"relationships"`
}

// GetTitle: Get a title for the chapter.
func (c *Chapter) GetTitle() string {
	return c.Attributes.Title
}

// GetChapterNum: Get the chapter's chapter number.
func (c *Chapter) GetChapterNum() string {
	if num := c.Attributes.Chapter; num != nil {
		return *num
	}
	return "-"
}

// ChapterAttributes: Attributes for a chapter.
type ChapterAttributes struct {
	Title              string  `json:"title"`
	Volume             *string `json:"volume"`
	Chapter            *string `json:"chapter"`
	TranslatedLanguage string  `json:"translatedLanguage"`
	Uploader           string  `json:"uploader"`
	ExternalURL        *string `json:"externalUrl"`
	IsUnavailable      bool    `json:"isUnavailable"`
	Pages              int     `json:"pages"`
	Version            int     `json:"version"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
	PublishAt          string  `json:"publishAt"`
	ReadableAt         string  `json:"readableAt"`
}

// Get: Get chapter by chapter id.
//
// https://api.mangadex.org/docs/redoc.html#tag/Chapter/operation/get-chapter-id
func (s *ChapterService) Get(id string, params url.Values) (chapter *Chapter, err error) {
	return s.GetContext(context.Background(), id, params)
}

// GetContext: Same as Get, with a context.
func (s *ChapterService) GetContext(ctx context.Context, id string, params url.Values) (chapter *Chapter, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = fmt.Sprintf(ChapterPath, id)
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("get chapter %q: %w", id, err)
	}
	err = json.Unmarshal(res.Data, &chapter)
	if err != nil {
		return nil, fmt.Errorf("get chapter %q: %w", id, err)
	}

	return chapter, nil
}

// List: Get chapter list.
//
// https://api.mangadex.org/docs/redoc.html#tag/Chapter/operation/get-chapter
func (s *ChapterService) List(params url.Values) (chapterList []*Chapter, err error) {
	return s.ListContext(context.Background(), params)
}

// ListContext: Same as List, with a context.
func (s *ChapterService) ListContext(ctx context.Context, params url.Values) (chapterList []*Chapter, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = ChapterListPath
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("list chapters: %w", err)
	}
	err = json.Unmarshal(res.Data, &chapterList)
	if err != nil {
		return nil, fmt.Errorf("list chapters: %w", err)
	}

	return chapterList, nil
}

// GetMangaChapters: Get a list of chapters for a manga by manga id.
//
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-manga-id-feed
//
// TODO: move this to manga.go?
func (s *ChapterService) GetMangaChapters(id string, params url.Values) (chapterList []*Chapter, err error) {
	return s.GetMangaChaptersContext(context.Background(), id, params)
}

// GetMangaChaptersContext: Same as GetMangaChapters, with a context.
func (s *ChapterService) GetMangaChaptersContext(ctx context.Context, id string, params url.Values) (chapterList []*Chapter, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = fmt.Sprintf(MangaChaptersPath, id)
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("list chapters for manga %q: %w", id, err)
	}
	err = json.Unmarshal(res.Data, &chapterList)
	if err != nil {
		return nil, fmt.Errorf("list chapters for manga %q: %w", id, err)
	}

	return chapterList, nil
}
