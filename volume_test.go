package mangodex

import (
	"net/url"
	"testing"
)

func TestVolume(t *testing.T) {
	// One Piece uuid
	id := "a1c7c817-4e59-43b7-9365-09675a149a6f"
	c := NewDexClient()
	params := url.Values{}
	params.Add("translatedLanguage[]", "en")
	resp, err := c.Volume.List(id, params)
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
	c := NewDexClient()
	params := url.Values{}
	params.Add("translatedLanguage[]", "en")
	resp, err := c.Volume.List(id, params)
	if err != nil {
		t.Fatal(err)
	}
	// Expected an empty map
	if len(resp) != 0 {
		t.Fatal(resp)
	}
	t.Log(resp)
}
