package mangodex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// ScanlationGroupListOptions: Options for the scanlation group list.
type ScanlationGroupListOptions struct {
	Limit           int       `json:"limit"`
	Offset          int       `json:"offset"`
	Ids             []string  `json:"ids"`
	Name            string    `json:"name"`
	FocusedLanguage string    `json:"focusedLanguage"`
	Includes        []string  `json:"includes"`
	Order           *GetOrder `json:"order,omitempty"`
}

// Get: Get scanlation group by scanlation group id.
// https://api.mangadex.org/docs/redoc.html#tag/ScanlationGroup/operation/get-group-id
func (s ScanlationGroupService) Get(id string, params url.Values) (*ScanlationGroup, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GroupGet, id)

	// Set query parameters
	u.RawQuery = params.Encode()

	res, err := s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var scanGroup ScanlationGroup
	err = json.Unmarshal(res.Data, &scanGroup)
	if err != nil {
		return nil, err
	}

	return &scanGroup, nil
}

// TODO: change this to use url.Values instead of custom option struct.
// This is the only method that has its custom options. Or add custom options to all other methods.
// List: Get scanlation group list.
// https://api.mangadex.org/docs/redoc.html#tag/ScanlationGroup/operation/get-search-group
func (s ScanlationGroupService) List(options *ScanlationGroupListOptions) ([]*ScanlationGroup, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(GroupList)

	// Set required query parameters
	q := u.Query()
	if options.FocusedLanguage != "" {
		q.Add("focusedLanguage", options.FocusedLanguage)
	}
	if options.Name != "" {
		q.Add("name", options.Name)
	}
	q.Add("limit", strconv.Itoa(options.Limit))
	q.Add("offset", strconv.Itoa(options.Offset))
	for _, i := range options.Ids {
		q.Add("ids[]", i)
	}
	for _, i := range options.Includes {
		q.Add("includes[]", i)
	}
	if options.Order != nil {
		// data, _ := json.Marshal(options.Order)
		// q.Add("order", string(data))
	}
	u.RawQuery = q.Encode()

	res, err := s.client.RequestAndDecode(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var scanGroups []*ScanlationGroup
	err = json.Unmarshal(res.Data, &scanGroups)
	if err != nil {
		return nil, err
	}

	return scanGroups, nil
}
