package command

import (
	"encoding/base64"
	"net/http"
	"net/url"

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

func (cli *StorageOSCli) WebsocketURLs() []*url.URL {

	urls := []*url.URL{}

	for _, h := range cli.GetHosts() {
		u, err := url.Parse(h)
		if err != nil {
			continue
		}
		u.Scheme = "ws"
		urls = append(urls, u)
	}
	return urls
}
