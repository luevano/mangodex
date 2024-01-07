package mangodex

import (
	"net/url"
	"testing"
)

func TestCover(t *testing.T) {
	c := NewDexClient()
	list, err := c.Manga.List(url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list[0].ID)

	resp, err := c.Cover.List([]string{list[0].ID}, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
