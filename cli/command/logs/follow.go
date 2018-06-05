package logs

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

const (

	// Endpoint path for logs
	logsPath = "/v1/logs"

	// Time allowed to read the next message from the peer.
	wait = 60 * time.Second

	// Send pings to peer with this period. Must be less than wait.
	pingPeriod = (wait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

func runFollow(storageosCli *command.StorageOSCli, opt logOptions) error {

	format := opt.format
	if len(format) == 0 {
		format = formatter.TableFormatKey
	}

	fmtCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewLogStreamFormat(format, opt.quiet),
	}

	c, err := storageosCli.WebsocketConn(logsPath)
	if err != nil || c == nil {
		return fmt.Errorf("Connection error: %v", err)
	}
	defer c.Close()

	c.SetPongHandler(func(string) error { return c.SetReadDeadline(time.Now().Add(wait)) })
	c.SetReadLimit(maxMessageSize)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Incoming message handler
	go func() {
		defer c.Close()
		defer close(done)
		for {
			if err := c.SetReadDeadline(time.Now().Add(wait)); err != nil {
				log.Error("Failed to SetReadDeadline:", err)
				return
			}
			_, message, err := c.ReadMessage()
			if err != nil {
				// Read errors are permanent, must close connection
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return
				}
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
					fmt.Println("Server closed connection")
					return
				}
				fmt.Printf("error: %v\n", err)
				return
			}

			// If node is specified, filter the logs as per nodes.
			if len(opt.nodes) > 0 {
				if skipNodeLog(opt.nodes, message) {
					continue
				}
			}

			formatter.LogStreamWrite(fmtCtx, message)
		}
	}()

	// Ticker for keepalives
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	// Write/Interrupt handler
	for {
		select {
		case <-interrupt:
			// Send a close frame and wait for the server to close the connection
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return nil
		case <-ticker.C:
			c.SetWriteDeadline(time.Now().Add(wait))
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("Keepalive failed: %v\n", err)
				return err
			}
		case <-done:
			return nil
		}
	}
}

// skipNodeLog tells if a message should be skipped based on the list of nodes.
func skipNodeLog(nodes []string, message []byte) bool {
	for _, node := range nodes {
		// Log message contains source node in the format "host":"storageos-1-80251".
		nodePattern := fmt.Sprintf("\"host\":\"%s\"", node)
		if bytes.Contains(message, []byte(nodePattern)) {
			return false
		}
	}

	return true
}
