package ws_client

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// this class attempts to emulate the twisted socket interface
// for usabilities sake and calls events on downstream nexuses
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

// base websocket service you can wrap in a struct so you don't need to reimplement
// empty event listeners
type BaseWebsocketClient struct {
	// See: https://git.io/JtPNQ, gorilla websockets do not suppor tconcurrent writers
	Mux            sync.Mutex
	Conn           *websocket.Conn
	BaseSharedSeed *beacon.SharedSeed
}

// do nothing
func (w *BaseWebsocketClient) OnConnecting() {}

func (w *BaseWebsocketClient) getConnection() *websocket.Conn {
	return w.Conn
}

func (w *BaseWebsocketClient) OnConnect(conn *websocket.Conn, r *http.Response) {
	w.Conn = conn
}
func (w *BaseWebsocketClient) OnOpen()                                        {}
func (w *BaseWebsocketClient) OnWsMessage(payload []byte, isBinary bool)      {}
func (w *BaseWebsocketClient) OnClose(wasClean bool, code int, reason string) {}
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

func (w *BaseWebsocketClient) SendBin(msg []byte) (err error) {
	if w.getConnection() == nil {
		return errors.New("not currently connected")
	}
	w.Mux.Lock()
	err = w.Conn.WriteMessage(websocket.BinaryMessage, msg)
	w.Mux.Unlock()
	return err
}

func (w *BaseWebsocketClient) SharedSeed() *beacon.SharedSeed {
	return w.BaseSharedSeed
}

var _ WebsocketClientProtocol = &BaseWebsocketClient{}

func NewBaseWebsocketClient() *BaseWebsocketClient {
	return &BaseWebsocketClient{}
}
