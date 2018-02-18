package command

import (
	"net/url"
	"testing"
)

func TestWebsocketURLs(t *testing.T) {
	url1, _ := url.Parse("ws://127.0.0.1")
	url2, _ := url.Parse("ws://storageos.net")

	testcases := []struct {
		name     string
		cli      StorageOSCli
		wantURLs []*url.URL
	}{
		{
			name: "numeric and name hosts",
			cli: StorageOSCli{
				hosts: []string{"127.0.0.1", "storageos.net"},
			},
			wantURLs: []*url.URL{url1, url2},
		},
		{
			name:     "no host",
			cli:      StorageOSCli{},
			wantURLs: []*url.URL{},
		},
		{
			name: "host with existing scheme",
			cli: StorageOSCli{
				hosts: []string{"https://storageos.net"},
			},
			wantURLs: []*url.URL{url2},
		},
		{
			name: "skip invalid host",
			cli: StorageOSCli{
				hosts: []string{":::....", "storageos.net"},
			},
			wantURLs: []*url.URL{url2},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotURLs := tc.cli.WebsocketURLs()
			if !urlsEq(gotURLs, tc.wantURLs) {
				t.Errorf("unexpected websocket URLs:\n\t(GOT) %v\n\t(WNT) %v", gotURLs, tc.wantURLs)
			}
		})
	}
}

// urlsEq compares slices of URLs.
func urlsEq(a, b []*url.URL) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].String() != b[i].String() {
			return false
		}
	}

	return true
}
