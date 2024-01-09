package mangodex

import (
	"net/url"
	"testing"
)

func TestCover(t *testing.T) {
	c := NewDexClient()
	params := url.Values{}
	params.Add("manga[]", "eadc095d-e672-4136-98d0-41a98161ad0e")

	resp, err := c.Cover.List(params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
