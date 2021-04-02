package server

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketServerProtocol attempts to emulate the twisted socket interface
// for usabilities sake and calls events on downstream nexuses
type WebSocketServerProtocol interface {
	// Context gets the context.Context object
	Context() context.Context
	// Cancel gets the context cancellation function
	Cancel() context.CancelFunc
	// OnConnecting is called as the websocket connects
	OnConnecting(r *http.Request)
	// OnConnect is called after the websocket conntects
	OnConnect(r *http.Request)
	// OnOpen is called after a message is opened
	OnOpen()
	// OnWsMessage processes a websocket message
	OnWsMessage(payload []byte, isBinary bool)
	// OnClose is called after the websocket service is closed
	OnClose(wasClean bool, code int, reason string)
	// ServeHTTP serves an http request
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// Upgrader is the websocket.Upgrader used to process ws requests
	Upgrader() websocket.Upgrader
	// SendMessage sends a message to the client
	SendMessage(msg []byte) error
	// Hub gets the Hub used to process websocket requests
	Hub() *Hub
}

// BaseWebsocketService you can wrap in a struct so you don't need to reimplement
// empty event listeners
type BaseWebsocketService struct {
	upgrader websocket.Upgrader
	hub      *Hub
	ctx      context.Context
	cancel   context.CancelFunc
}

// SendMessage sends a message to the client
func (w BaseWebsocketService) SendMessage(msg []byte) error {
	w.hub.broadcast <- msg
	return nil
}

// Hub gets the Hub used to process websocket requests
func (w BaseWebsocketService) Hub() *Hub {
	return w.hub
}

// Context gets the context.Context object
func (w BaseWebsocketService) Context() context.Context {
	return w.ctx
}

// OnConnecting is called as the websocket connects
// implemented in a sub-struct
func (w BaseWebsocketService) OnConnecting(r *http.Request) {}

// OnConnect is called after the websocket conntects
// implemented in a sub-struct
func (w BaseWebsocketService) OnConnect(r *http.Request) {}

// OnOpen is called after a message is opened
// implemented in a sub-struct
func (w BaseWebsocketService) OnOpen() {}

// OnWsMessage processes a websocket message
// implemented in a sub-struct
func (w BaseWebsocketService) OnWsMessage(payload []byte, isBinary bool) {}

// OnClose is called after the websocket service is closed
// implemented in a sub-struct
func (w BaseWebsocketService) OnClose(wasClean bool, code int, reason string) {}

// ServeHTTP serves an http request
// implemented in a sub-struct
func (w BaseWebsocketService) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	panic("method must be implemented on child")
}

// Upgrader is the websocket.Upgrader used to process ws requests
// implemented in a sub-struct
func (w BaseWebsocketService) Upgrader() websocket.Upgrader {
	return w.upgrader
}

// Cancel gets the context cancellation function
// implemented in a sub-struct
func (w BaseWebsocketService) Cancel() context.CancelFunc {
	return w.cancel
}

var _ WebSocketServerProtocol = &BaseWebsocketService{}

// NewBaseWebsocketService creates a BaseWebsocketService
func NewBaseWebsocketService() BaseWebsocketService {
	ctx, cancel := context.WithCancel(context.Background())
	bss := BaseWebsocketService{
		upgrader: websocket.Upgrader{
			// TODO
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		hub:    NewHub(),
		ctx:    ctx,
		cancel: cancel,
	}
	go bss.hub.run()
	return bss
}
