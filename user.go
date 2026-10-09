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
	return s.GetContext(context.Background(), id)
}

// GetContext: Same as Get, with a context.
func (s *UserService) GetContext(ctx context.Context, id string) (user *User, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = fmt.Sprintf(GetUserPath, id)

	var res UserResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("get user %q: %w", id, err)
	}
	err = json.Unmarshal(res.Data, &user)
	if err != nil {
		return nil, fmt.Errorf("get user %q: %w", id, err)
	}

	return user, nil
}
