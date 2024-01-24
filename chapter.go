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
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(ChapterPath, id)
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Data, &chapter)
	if err != nil {
		return nil, err
	}

	return chapter, nil
}

// List: Get chapter list.
//
// https://api.mangadex.org/docs/redoc.html#tag/Chapter/operation/get-chapter
func (s *ChapterService) List(params url.Values) (chapterList []*Chapter, err error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = ChapterListPath
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Data, &chapterList)
	if err != nil {
		return nil, err
	}

	return chapterList, nil
}

// GetMangaChapters: Get a list of chapters for a manga by manga id.
//
// https://api.mangadex.org/docs/redoc.html#tag/Manga/operation/get-manga-id-feed
//
// TODO: move this to manga.go?
func (s *ChapterService) GetMangaChapters(id string, params url.Values) (chapterList []*Chapter, err error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaChaptersPath, id)
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Data, &chapterList)
	if err != nil {
		return nil, err
	}

	return chapterList, nil
}

// TODO: update viable methods later. Most of this is either deprecated
// or the API changed drastically (due to auth being different).
// The code is heavily outdated.
/*
// ChapterReadMarkers: A response for getting a list of read chapters.
type ChapterReadMarkers struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

func (rmr *ChapterReadMarkers) GetResult() string {
	return rmr.Result
}

// GetReadMangaChapters: Get list of Chapter IDs that are marked as read for a specified manga ID.
// https://api.mangadex.org/docs.html#operation/get-manga-chapter-readmarkers
func (s *ChapterService) GetReadMangaChapters(id string) (*ChapterReadMarkers, error) {
	return s.GetReadMangaChaptersContext(context.Background(), id)
}

// GetReadMangaChaptersContext: GetReadMangaChapters with custom context.
func (s *ChapterService) GetReadMangaChaptersContext(ctx context.Context, id string) (*ChapterReadMarkers, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaReadMarkersPath, id)

	var rmr ChapterReadMarkers
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &rmr)
	return &rmr, err
}

// SetReadUnreadMangaChapters: Set read/unread manga chapters.
func (s *ChapterService) SetReadUnreadMangaChapters(id string, read, unRead []string) (*Response, error) {
	return s.SetReadUnreadMangaChaptersContext(context.Background(), id, read, unRead)
}

// SetReadUnreadMangaChaptersContext: SetReadUnreadMangaChapters with custom context.
func (s *ChapterService) SetReadUnreadMangaChaptersContext(ctx context.Context, id string, read, unRead []string) (*Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaReadMarkersPath, id)

	// Set request body.
	req := map[string][]string{
		"chapterIdsRead":   read,
		"chapterIdsUnread": unRead,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r Response
	err = s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), bytes.NewBuffer(rBytes), &r)
	return &r, err
}
*/
