package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GroupList = "/group"
	GroupGet  = "/group/%s"
)

// ScanlationGroupService: Provides scanlation group services provided by the API.
type ScanlationGroupService service

// ScanlationGroup: Struct containing information on a scanlation group.
type ScanlationGroup struct {
	Id            string                     `json:"id"`
	Type          RelationshipType           `json:"type"`
	Attributes    *ScanlationGroupAttributes `json:"attributes"`
	Relationships []*Relationship            `json:"relationships"`
}

// ScanlationGroupAttributes: Attributes for a scanlation group
type ScanlationGroupAttributes struct {
	Name            string           `json:"name"`
	AltNames        LocalisedStrings `json:"altNames"`
	Website         *string          `json:"website"`
	IRCServer       *string          `json:"ircServer"`
	Discord         *string          `json:"discord"`
	ContactEmail    *string          `json:"contactEmail"`
	Description     *string          `json:"description"`
	Twitter         *string          `json:"twitter"`
	FocusedLanguage []string         `json:"focusedLanguage"`
	Locked          bool             `json:"locked"`
	Official        bool             `json:"official"`
	Inactive        bool             `json:"inactive"`
	PublishDelay    string           `json:"publishDelay"`
	Version         int              `json:"version"`
	CreatedAt       string           `json:"createdAt"`
	UpdatedAt       string           `json:"updatedAt"`
}

// Get: Get scanlation group by scanlation group id.
//
// https://api.mangadex.org/docs/redoc.html#tag/ScanlationGroup/operation/get-group-id
func (s *ScanlationGroupService) Get(id string, params url.Values) (group *ScanlationGroup, err error) {
	return s.GetContext(context.Background(), id, params)
}

// GetContext: Same as Get, with a context.
func (s *ScanlationGroupService) GetContext(ctx context.Context, id string, params url.Values) (group *ScanlationGroup, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = fmt.Sprintf(GroupGet, id)
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("get scanlation group %q: %w", id, err)
	}
	err = json.Unmarshal(res.Data, &group)
	if err != nil {
		return nil, fmt.Errorf("get scanlation group %q: %w", id, err)
	}

	return group, nil
}

// List: Get scanlation group list.
//
// https://api.mangadex.org/docs/redoc.html#tag/ScanlationGroup/operation/get-search-group
func (s *ScanlationGroupService) List(params url.Values) (groupList []*ScanlationGroup, err error) {
	return s.ListContext(context.Background(), params)
}

// ListContext: Same as List, with a context.
func (s *ScanlationGroupService) ListContext(ctx context.Context, params url.Values) (groupList []*ScanlationGroup, err error) {
	u, _ := url.Parse(s.client.baseAPI)
	u.Path = GroupList
	u.RawQuery = params.Encode()

	var res DexResponse
	err = s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("list scanlation groups: %w", err)
	}
	err = json.Unmarshal(res.Data, &groupList)
	if err != nil {
		return nil, fmt.Errorf("list scanlation groups: %w", err)
	}

	return groupList, nil
}
