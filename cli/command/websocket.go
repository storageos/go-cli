package command

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

func (cli *StorageOSCli) WebsocketConn(path string) (*websocket.Conn, error) {
	authHeader := base64.StdEncoding.EncodeToString([]byte(cli.GetUsername() + ":" + cli.GetPassword()))

	var c *websocket.Conn
	var err error

	for _, u := range cli.WebsocketURLs() {
		u.Path = path
		c, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": {authHeader}})
		if err == nil {
			break
		}
	}

	return c, err
}

// WebsocketURLs creates websocket URL of all the hosts and returns a slice of
// the URLs.
func (cli *StorageOSCli) WebsocketURLs() []*url.URL {
	urls := []*url.URL{}
	hosts := strings.Split(cli.hosts, ",")
	for _, h := range hosts {
		u, err := url.Parse(h)
		if err != nil {
			continue
		}
		u.Scheme = "ws"
		urls = append(urls, u)
	}
	return urls
}
