package relay

import (
	"log"
	"time"

	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/relay"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// Relay relays messages to/from terminus.
type Relay struct {
	// config is the terminus config used for starting the relay.
	// technically, we could use config.RelayConfig here, as nothing else is (should be) used
	config *config.Config
	// websocketLayer is the websocket layer used for communications with browsers
	websocketLayer *websocket.IncomingWebsocketLayer
	// rendezvousLayer is responsible for rendezvousing with different nexuses
	rendezvousLayer *rendezvous.IncomingRendezvousLayer
	// relayLayer is responsible for relaying messages to their peer nexuses
	relayLayer *relay.Layer
}

// NewRelay creates a new relay from a config and starts a looping info call.
func NewRelay(config *config.Config) Relay {
	r := Relay{}
	r.config = config
	r.SetupWebsocketLayer()
	r.SetupRendezvousLayer(r.websocketLayer)
	r.SetupRelayLayer(r.rendezvousLayer)
	go r.OutputInfo()
	return r
}

// OnStackEvent is used for any events that bubble up to Relay.
func (r *Relay) OnStackEvent(layerName string, nexus nexusHelper.Nexus, event string) {
	// do nothing for now
}

// OutputInfo gets nexus count on a relay.
func (r *Relay) OutputInfo() {
	for {
		<-time.After(2 * time.Second)
		log.Println(r.rendezvousLayer.ToString())
	}
}

// SetupWebsocketLayer sets up a websocket.WebsocketLayer and registers it.
func (r *Relay) SetupWebsocketLayer() *Relay {
	r.websocketLayer = websocket.NewIncomingWebsocketLayer(r.config)
	r.websocketLayer.RegisterLayerEvent(r.OnStackEvent, message.IncomingWebsocket)
	return r
}

// SetupRendezvousLayer sets up a rendezvous.Rendezvous layer and registers it.
func (r *Relay) SetupRendezvousLayer(belowLayer layer.Base) *Relay {
	r.rendezvousLayer = rendezvous.NewIncomingRendezvousLayer()
	r.rendezvousLayer.RegisterAboveLayer(belowLayer)
	r.rendezvousLayer.RegisterLayerEvent(r.OnStackEvent, message.IncomingRendezvous)
	return r
}

// SetupRelayLayer sets up a relay.RelayLayer and registers it.
func (r *Relay) SetupRelayLayer(rendezvousLayer *rendezvous.IncomingRendezvousLayer) *Relay {
	r.relayLayer = relay.NewRelayLayer(rendezvousLayer)
	r.relayLayer.SetOnAnnounce(rendezvousLayer.OnAnnounce)
	r.relayLayer.SetOnRevoke(rendezvousLayer.OnRevoke)
	return r
}

// RunApp listens on the relay url.
func (r *Relay) RunApp() error {
	return r.websocketLayer.Listen(r.config.GetRelayURL(), nil)
}
