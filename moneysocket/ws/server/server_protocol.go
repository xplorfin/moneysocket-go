package server

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

// this class attempts to emulate the twisted socket interface
// for usabilities sake and calls events on downstream nexuses
type WebSocketServerProtocol interface {
	// context
	Context() context.Context
	// cancel func
	Cancel() context.CancelFunc
	// calls when a new connection is made
	OnConnecting(r *http.Request)
	// the client has connected
	OnConnect(r *http.Request)
	// connection is open
	OnOpen()
	// receive a message, maps to OnMessage() in python version
	OnWsMessage(payload []byte, isBinary bool)
	// called when connection isclosed
	OnClose(wasClean bool, code int, reason string)
	// serve http interface. Must be implemented on child
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// get the upgrader
	Upgrader() websocket.Upgrader

	SendMessage(msg []byte) error
	// get hub
	Hub() *Hub
}

// base websocket service you can wrap in a struct so you don't need to reimplement
// empty event listeners
type BaseWebsocketService struct {
	upgrader websocket.Upgrader
	hub      *Hub
	ctx      context.Context
	cancel   context.CancelFunc
}

func (w BaseWebsocketService) SendMessage(msg []byte) error {
	w.hub.broadcast <- msg
	return nil
}
func (w BaseWebsocketService) Hub() *Hub {
	return w.hub
}

func (w BaseWebsocketService) Context() context.Context {
	return w.ctx

}
func (w BaseWebsocketService) OnConnecting(r *http.Request)                   {}
func (w BaseWebsocketService) OnConnect(r *http.Request)                      {}
func (w BaseWebsocketService) OnOpen()                                        {}
func (w BaseWebsocketService) OnWsMessage(payload []byte, isBinary bool)      {}
func (w BaseWebsocketService) OnClose(wasClean bool, code int, reason string) {}
func (w BaseWebsocketService) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	panic("method must be implemented on child")
}
func (w BaseWebsocketService) Upgrader() websocket.Upgrader {
	return w.upgrader
}

func (w BaseWebsocketService) Cancel() context.CancelFunc {
	return w.cancel

}

var _ WebSocketServerProtocol = &BaseWebsocketService{}

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
