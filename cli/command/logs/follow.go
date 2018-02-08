package logs

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

const logsPath = "/v1/logs"

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
		return fmt.Errorf("Could not create connection to stream logs")
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Incoming message handler
	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				// Read errors are permanent, must close connection
				return
			}
			formatter.LogStreamWrite(fmtCtx, message)
		}
	}()

	// Interrupt handler
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
		case <-done:
			return nil
		}
	}
}
