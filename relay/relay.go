package relay

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/relay"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

type Relay struct {
	config          *config.Config
	websocketLayer  *websocket.IncomingWebsocketLayer
	rendezvousLayer *rendezvous.IncomingRendezvousLayer
	relayLayer      *relay.RelayLayer
}

func NewRelay(config *config.Config) Relay {
	r := Relay{}
	r.config = config
	r.SetupWebsocketLayer()
	r.SetupRendezvousLayer(r.websocketLayer)
	r.SetupRelayLayer(r.rendezvousLayer)
	return r
}

func (r *Relay) OnStackEvent(layerName string, nexus nexusHelper.Nexus, event string) {
	// do nothing for now
}

func (r *Relay) SetupWebsocketLayer() *Relay {
	r.websocketLayer = websocket.NewIncomingWebsocketLayer(r.config)
	r.websocketLayer.RegisterLayerEvent(r.OnStackEvent, message.IncomingWebsocket)
	return r
}

func (r *Relay) SetupRendezvousLayer(belowLayer layer.Layer) *Relay {
	r.rendezvousLayer = rendezvous.NewIncomingRendezvousLayer()
	r.rendezvousLayer.RegisterAboveLayer(belowLayer)
	r.rendezvousLayer.RegisterLayerEvent(r.OnStackEvent, message.IncomingRendezvous)
	return r
}

func (r *Relay) SetupRelayLayer(rendezvousLayer *rendezvous.IncomingRendezvousLayer) *Relay {
	r.relayLayer = relay.NewRelayLayer(rendezvousLayer)
	r.relayLayer.SetOnAnnounce(rendezvousLayer.OnAnnounce)
	r.relayLayer.SetOnRevoke(rendezvousLayer.OnRevoke)
	return r
}

func (r *Relay) RunApp() error {
	return r.websocketLayer.Listen(r.config.GetRelayUrl(), nil)
}
