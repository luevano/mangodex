package mangodex

import (
	"net/url"
	"strconv"
	"testing"
)

// TODO: refactor all the tests

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

//
// manga.go
//

func TestMangaGet(t *testing.T) {
	manga, err := client.Manga.Get("acc3ff9c-3494-4bdc-b474-96b24d0c160c", url.Values{})
	if err != nil && manga.Attributes.Title.GetLocalString("en") == "Ochite Oborete" {
		t.Errorf("Getting Manga from Group failed: %s\n", err.Error())
	}
}

func TestMangaList(t *testing.T) {
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

// scanlation_group.go

func TestGroupGet(t *testing.T) {
	group, err := client.ScanlationGroup.Get("71ade5cd-93cf-4397-a5cc-d5c6181d8697", url.Values{})
	if err != nil && group.Id == "71ade5cd-93cf-4397-a5cc-d5c6181d8697" {
		t.Errorf("Getting Group failed: %s\n", err.Error())
	}
}

func TestGroupList(t *testing.T) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(100))
	params.Set("offset", strconv.Itoa(0))
	params.Set("name", "laughing")
	_, err := client.ScanlationGroup.List(params)
	if err != nil {
		t.Errorf("Getting Group List failed: %s\n", err.Error())
	}
}

//
// chapter.go
//

func TestChapterGet(t *testing.T) {
	chapter, err := client.Chapter.Get("eadc095d-e672-4136-98d0-41a98161ad0e", url.Values{})
	if err != nil && chapter.GetTitle() == "Attachment" {
		t.Errorf("Getting Manga from Group failed: %s\n", err.Error())
	}
}

//
// cover.go
//

func TestCover(t *testing.T) {
	params := url.Values{}
	params.Add("manga[]", "eadc095d-e672-4136-98d0-41a98161ad0e")

	resp, err := client.Cover.List(params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

//
// volume.go
//

func TestVolume(t *testing.T) {
	// One Piece uuid
	id := "a1c7c817-4e59-43b7-9365-09675a149a6f"
	params := url.Values{}
	params.Add("translatedLanguage[]", "en")
	resp, err := client.Volume.List(id, params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

// TODO: change to a more reliable test, this can change at any time whenever some volumes/chapters are added
func TestVolumeEmpty(t *testing.T) {
	// Demon Slayer (not Kimetsu No Yaiba) uuid, at the time of writing this it doesn't have any chapters/volumes
	// https://mangadex.org/title/0acb51ef-3d71-4993-81a0-8cbcfb88fa9e/demon-slayer
	id := "0acb51ef-3d71-4993-81a0-8cbcfb88fa9e"
	params := url.Values{}
	params.Add("translatedLanguage[]", "en")
	resp, err := client.Volume.List(id, params)
	if err != nil {
		t.Fatal(err)
	}
	// Expected an empty map
	if len(resp) != 0 {
		t.Fatal(resp)
	}
	t.Log(resp)
}

//
// user.go
//

func TestUser(t *testing.T) {
	// Newtonius uuid
	id := "904b5ab6-7e00-4b7e-a6c6-3dda7860b69e"
	resp, err := client.User.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
