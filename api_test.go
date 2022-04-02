package mangodex

import (
	"net/url"
	"strconv"
	"testing"
)

var client = NewDexClient()

/*
func TestLogin(t *testing.T) {
	err := client.Auth.Login(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		t.Error("Login failed.")
	}
	fmt.Printf("%v\n", client)
}

func TestGetLoggedUser(t *testing.T) {
	user, err := client.User.GetLoggedUser()
	if err != nil {
		t.Error("Getting user failed.")
	}
	t.Log(user)
}
*/

func TestGetMangaList(t *testing.T) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(100))
	params.Set("offset", strconv.Itoa(0))
	// Include Author relationship
	params.Set("includes[]", string(RelationshipTypeAuthor))
	// If it is a search, then we add the search term.
	_, err := client.Manga.List(params)
	if err != nil {
		t.Errorf("Getting manga failed: %s\n", err.Error())
	}
}

func TestGroupList(t *testing.T) {
	_, err := client.ScanlationGroup.List(&ScanlationGroupListOptions{
		Limit:           100,
		Offset:          0,
		Ids:             nil,
		Name:            "laughing",
		FocusedLanguage: "",
		Includes:        nil,
		Order:           nil,
	})
	if err != nil {
		t.Errorf("Getting Group List failed: %s\n", err.Error())
	}
}

func TestGroupGet(t *testing.T) {
	_, err := client.ScanlationGroup.Get("71ade5cd-93cf-4397-a5cc-d5c6181d8697")
	if err != nil {
		t.Errorf("Getting Group failed: %s\n", err.Error())
	}
}

func TestMangaGroupList(t *testing.T) {
	list, err := client.Manga.List(url.Values{
		"limit":      {"100"},
		"group":      {"71ade5cd-93cf-4397-a5cc-d5c6181d8697"},
		"includes[]": {string(RelationshipTypeCoverArt), string(RelationshipTypeArtist), string(RelationshipTypeAuthor)},
	})
	if err != nil && len(list) > 0 {
		t.Errorf("Getting Manga from Group failed: %s\n", err.Error())
	}
}

func TestMangaGet(t *testing.T) {
	manga, err := client.Manga.Get("acc3ff9c-3494-4bdc-b474-96b24d0c160c")
	if err != nil && manga.Attributes.Title.GetLocalString("en") == "Ochite Oborete" {
		t.Errorf("Getting Manga from Group failed: %s\n", err.Error())
	}
}
