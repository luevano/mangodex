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
