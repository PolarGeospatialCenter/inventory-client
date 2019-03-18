package client

import (
	"net/url"
	"testing"
)

func TestInventoryUrl(t *testing.T) {
	baseUrl, _ := url.Parse("https://inventory.api.local/v0")
	inv := NewInventoryApi(baseUrl)

	nodeUrl := inv.Url("/node")
	if nodeUrl != "https://inventory.api.local/v0/node" {
		t.Errorf("got wrong node url: %s", nodeUrl)
	}

}
