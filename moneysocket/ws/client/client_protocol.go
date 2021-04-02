package client

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// WebsocketClientProtocol class attempts to emulate the twisted socket interface
// for usabilities sake and calls events on downstream nexuses.
type WebsocketClientProtocol interface {
	// calls when a new connection is made
	OnConnecting()
	// the client has connected
	OnConnect(conn *websocket.Conn, r *http.Response)
	// connection is open
	OnOpen()
	// handle incoming messages
	OnWsMessage(payload []byte, isBinary bool)
	// called when connection isclosed
	OnClose(wasClean bool, code int, reason string)
	// get the websocket connection object
	getConnection() *websocket.Conn
	// send a message
	Send(msg base.MoneysocketMessage) error
	// send a binary-encoded message
	SendBin(msg []byte) error
	// get the shared seed
	SharedSeed() *beacon.SharedSeed
}

// BaseWebsocketClient base websocket service you can wrap in a struct so you don't need to reimplement
// empty event listeners. It  also provides a canonical way to send messages.
type BaseWebsocketClient struct {
	// See: https://git.io/JtPNQ, gorilla websockets do not suppor tconcurrent writers
	Mux            sync.Mutex
	Conn           *websocket.Conn
	BaseSharedSeed *beacon.SharedSeed
}

// OnConnecting is called when connection is established.
func (w *BaseWebsocketClient) OnConnecting() {}

// getConnection gets the connection object.
func (w *BaseWebsocketClient) getConnection() *websocket.Conn {
	return w.Conn
}

// OnConnect called on connect, sets the connection object.
func (w *BaseWebsocketClient) OnConnect(conn *websocket.Conn, r *http.Response) {
	w.Conn = conn
}

// OnOpen called when the connection is opened.
func (w *BaseWebsocketClient) OnOpen() {}

// OnWsMessage called when a ws message is received.
func (w *BaseWebsocketClient) OnWsMessage(payload []byte, isBinary bool) {}

// OnClose called when a connection is closed.
func (w *BaseWebsocketClient) OnClose(wasClean bool, code int, reason string) {}

// Send a websocket message.
func (w *BaseWebsocketClient) Send(msg base.MoneysocketMessage) error {
	if w.getConnection() == nil {
		return errors.New("not currently connected")
	}
	res, err := message.WireEncode(msg, w.SharedSeed())
	if err != nil {
		return err
	}
	return w.SendBin(res)
}

// SendBin sends a binary message.
func (w *BaseWebsocketClient) SendBin(msg []byte) (err error) {
	if w.getConnection() == nil {
		return errors.New("not currently connected")
	}
	w.Mux.Lock()
	err = w.Conn.WriteMessage(websocket.BinaryMessage, msg)
	w.Mux.Unlock()
	return err
}

// SharedSeed gets the shared seed.
func (w *BaseWebsocketClient) SharedSeed() *beacon.SharedSeed {
	return w.BaseSharedSeed
}

// make sure the client binds to the base websocket client.
var _ WebsocketClientProtocol = &BaseWebsocketClient{}

// NewBaseWebsocketClient creates a websocket client with default options.
func NewBaseWebsocketClient() *BaseWebsocketClient {
	return &BaseWebsocketClient{}
}
