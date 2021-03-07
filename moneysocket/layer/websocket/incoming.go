package websocket

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/ws/ws_server"
)

type IncomingWebsocketLayer struct {
	layer.BaseLayer
	Config         *config.Config
	IncomingSocket *websocket.IncomingSocket
	WebsocketNexus *websocket.WebsocketNexus
}

func (i *IncomingWebsocketLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(i.OnAnnounce)
	belowLayer.SetOnRevoke(i.OnRevoke)
}

func NewIncomingWebsocketLayer(config *config.Config) *IncomingWebsocketLayer {
	wn := websocket.NewIncomingSocket()
	is := IncomingWebsocketLayer{
		BaseLayer:      layer.NewBaseLayer(),
		Config:         config,
		IncomingSocket: &wn,
	}
	// set factory ms protocol layer to the current layer
	is.IncomingSocket.FactoryMsProtocolLayer = &is
	return &is
}

func (i *IncomingWebsocketLayer) AnnounceNexus(belowNexus nexusHelper.Nexus) {
	websocketNexus := websocket.NewWebsocketNexus(belowNexus, i)
	i.WebsocketNexus = &websocketNexus

	i.TrackNexus(i.WebsocketNexus, belowNexus)
	i.TrackNexusAnnounced(i.WebsocketNexus)
	i.SendLayerEvent(i.WebsocketNexus, message.NexusAnnounced)
	if i.OnAnnounce != nil {
		i.OnAnnounce(i.WebsocketNexus)
	}
}

func (i *IncomingWebsocketLayer) StopListening() {
	panic("method not yet implemented")
}

func (i *IncomingWebsocketLayer) Listen(wsUrl string, tlsInfo *ws_server.TlsInfo) (err error) {
	i.IncomingSocket.FactoryMsProtocolLayer = i
	return ws_server.Listen(wsUrl, tlsInfo, i.IncomingSocket.ServeHTTP)
}

var _ layer.Layer = &IncomingWebsocketLayer{}
