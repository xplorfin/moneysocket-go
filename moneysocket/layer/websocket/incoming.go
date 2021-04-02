package websocket

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/ws/server"
)

// IncomingWebsocketLayer is handles websockets.
type IncomingWebsocketLayer struct {
	layer.BaseLayer
	Config *config.Config
	// IncomingSocket is a simulated socket
	IncomingSocket *websocket.IncomingSocket
	// WebsocketNexus handles the socket
	WebsocketNexus *websocket.Nexus
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer.
func (i *IncomingWebsocketLayer) RegisterAboveLayer(belowLayer layer.Base) {
	belowLayer.SetOnAnnounce(i.OnAnnounce)
	belowLayer.SetOnRevoke(i.OnRevoke)
}

// NewIncomingWebsocketLayer creates a new incoming websocket layer.
func NewIncomingWebsocketLayer(config *config.Config) *IncomingWebsocketLayer {
	wn := websocket.NewIncomingSocket()
	is := IncomingWebsocketLayer{
		BaseLayer:      layer.NewBaseLayer(),
		Config:         config,
		IncomingSocket: wn,
	}
	// set factory ms protocol layer to the current layer
	is.IncomingSocket.FactoryMsProtocolLayer = &is
	return &is
}

// AnnounceNexus creates a new websocket.WebsocketNexus and registers it.
func (i *IncomingWebsocketLayer) AnnounceNexus(belowNexus nexusHelper.Nexus) {
	websocketNexus := websocket.NewWebsocketNexus(belowNexus, i)
	i.WebsocketNexus = websocketNexus

	i.TrackNexus(i.WebsocketNexus, belowNexus)
	i.TrackNexusAnnounced(i.WebsocketNexus)
	i.SendLayerEvent(i.WebsocketNexus, message.NexusAnnounced)
	if i.OnAnnounce != nil {
		i.OnAnnounce(i.WebsocketNexus)
	}
}

// StopListening stops listening on the incoming websocket layer.
func (i *IncomingWebsocketLayer) StopListening() {
	panic("method not yet implemented")
}

// Listen listens at a given wsUrl.
func (i *IncomingWebsocketLayer) Listen(wsURL string, tlsInfo *server.TLSInfo) (err error) {
	i.IncomingSocket.FactoryMsProtocolLayer = i
	return server.Listen(wsURL, tlsInfo, i.IncomingSocket.ServeHTTP)
}

var _ layer.Base = &IncomingWebsocketLayer{}
