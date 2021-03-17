package ws_client

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// code is setnt to this when code is unknown
const UnknownStatusCode = -1

// TODO make this do things
func NewWsClient(p WebsocketClientProtocol, wsUrl string) {
	p.OnConnecting()

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}

	c, res, err := dialer.Dial(wsUrl, nil)
	if err != nil {
		statusCode := UnknownStatusCode
		// handle cases where no response code is returned
		if res != nil {
			statusCode = res.StatusCode
		}
		p.OnClose(false, statusCode, err.Error())
		return
	}

	p.OnConnect(c, res)
	p.OnOpen()

	done := make(chan struct{})

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		defer close(done)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				p.OnClose(false, UnknownStatusCode, err.Error())
				return
			}
			p.OnWsMessage(message, mt == websocket.BinaryMessage)
		}
	}()
}
