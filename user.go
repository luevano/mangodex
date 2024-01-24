package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GetUserPath = "/user/%s"
	// GetUserFollowedMangaListPath = "/user/follows/manga"
	// GetLoggedUserPath            = "/user/me"
)

// UserResponse: User response type, this differs from the common DexResponse type.
type UserResponse struct {
	Result   string          `json:"result"`
	Response string          `json:"response"`
	Data     json.RawMessage `json:"data"`
}

// UserService: Provides user services provided by the API.
type UserService service

// User: Info on a MangaDex user.
type User struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Attributes    UserAttributes `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
}

// UserAttributes: Attributes of a User.
type UserAttributes struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Version  int      `json:"version"`
}

// Get: Get user by id.
//
// https://api.mangadex.org/docs/redoc.html#tag/User/operation/get-user-id
func (s *UserService) Get(id string) (user *User, err error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GetUserPath, id)

	var res UserResponse
	err = s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Data, &user)
	if err != nil {
		return nil, err
	}

	return user, err
}

// TODO: enable once Auth service is fixed.
/*
// GetUserFollowedMangaList: Get list of followed manga.
//
// https://api.mangadex.org/docs/redoc.html#tag/Follows/operation/get-user-follows-manga
func (s *UserService) GetUserFollowedMangaList(params url.Values) ([]*Manga, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetUserFollowedMangaListPath
	u.RawQuery = params.Encode()

	res, err := s.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var mangaList []*Manga
	err = json.Unmarshal(res.Data, &mangaList)
	if err != nil {
		return nil, err
	}

	return mangaList, err
}

// GetLoggedUser: Get logged user.
//
// https://api.mangadex.org/docs/redoc.html#tag/User/operation/get-user-me
func (s *UserService) GetLoggedUser() (*User, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetLoggedUserPath

	res, err := s.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var user User
	err = json.Unmarshal(res.Data, &user)
	if err != nil {
		return nil, err
	}

	return &user, err
}
*/
