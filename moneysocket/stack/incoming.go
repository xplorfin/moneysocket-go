package stack

import (
	"fmt"
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/local"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/relay"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// IncomingStack handles incoming connections
type IncomingStack struct {
	layer.BaseLayer
	config          *config.Config
	localLayer      *local.IncomingLocalLayer
	websocketLayer  *websocket.IncomingWebsocketLayer
	rendezvousLayer *rendezvous.IncomingRendezvousLayer
	relayLayer      *relay.Layer
}

// NewIncomingStack creates a new incoming stack
func NewIncomingStack(config *config.Config, outgoingLocalLayer *local.OutgoingLocalLayer) *IncomingStack {
	is := IncomingStack{
		BaseLayer: layer.NewBaseLayer(),
		config:    config,
	}
	is.SetupLocalLayer(outgoingLocalLayer)
	is.SetupWebsocketLayer()
	is.SetupRendezvousLayer(is.websocketLayer, is.localLayer)
	is.SetupRelayLayer(is.rendezvousLayer)
	return &is
}

// SetupRelayLayer sets up the relay layer
func (i *IncomingStack) SetupRelayLayer(rendezvousLayer *rendezvous.IncomingRendezvousLayer) {
	i.relayLayer = relay.NewRelayLayer(rendezvousLayer)
	i.relayLayer.RegisterAboveLayer(rendezvousLayer)
	i.relayLayer.RegisterLayerEvent(i.SendStackEvent, message.Relay)
}

// SetupRendezvousLayer sets up the rendezvous layer
func (i *IncomingStack) SetupRendezvousLayer(belowLayer1 layer.Base, belowLayer2 layer.Base) {
	i.rendezvousLayer = rendezvous.NewIncomingRendezvousLayer()
	i.rendezvousLayer.RegisterAboveLayer(belowLayer1)
	i.rendezvousLayer.RegisterAboveLayer(belowLayer2)
	i.rendezvousLayer.RegisterLayerEvent(i.SendStackEvent, message.IncomingRendezvous)
}

// SetupWebsocketLayer sets up the websocket layer
func (i *IncomingStack) SetupWebsocketLayer() {
	i.websocketLayer = websocket.NewIncomingWebsocketLayer(i.config)
	i.websocketLayer.RegisterLayerEvent(i.SendStackEvent, message.IncomingWebsocket)
}

// SetupLocalLayer sets up the outgoing local layer
func (i *IncomingStack) SetupLocalLayer(outgoingLocalLayer *local.OutgoingLocalLayer) {
	i.localLayer = local.NewIncomingLocalLayer()
	i.localLayer.RegisterLayerEvent(i.SendStackEvent, message.IncomingLocal)
	outgoingLocalLayer.SetIncomingLayer(i.localLayer)
}

// SendStackEvent does nothing in this context
func (i *IncomingStack) SendStackEvent(layerName string, nexus nexusHelper.Nexus, event string) {
	// do nothing
}

// RegisterAboveLayer does nothing. It is here to satisfy the interface since IncomingStack must
// act as a layer
func (i *IncomingStack) RegisterAboveLayer(belowLayer layer.Base) {
	// do nothing
}

// AnnounceNexus does nothing. It is here to satisfy the interface since IncomingStack must
// act as a layer
func (i *IncomingStack) AnnounceNexus(belowNexus nexusHelper.Nexus) {
	log.Println("announced from below")
}

// RevokeNexus does nothing. It is here to satisfy the interface since IncomingStack
// acts as a layer
func (i *IncomingStack) RevokeNexus(nexus nexusHelper.Nexus) {
	log.Println("revoked from below")
}

// GetListenLocations returns a websocket location object from config
func (i *IncomingStack) GetListenLocations() []location.Location {
	return []location.Location{location.NewWebsocketLocationPort(i.config.GetExternalHost(), i.config.GetUseTLS(), i.config.GetExternalPort())}
}

// GetListenURL gets a listen url
func (i *IncomingStack) GetListenURL() string {
	schema := "ws"
	if i.config.GetUseTLS() {
		schema = "wss"
	}
	return fmt.Sprintf("%s://%s:%d", schema, i.config.ListenConfig.BindHost, i.config.ListenConfig.BindPort)
}

// Listen listens on a given port
// TODO implement tls config
func (i *IncomingStack) Listen() error {
	return i.websocketLayer.Listen(i.GetListenURL(), nil)
}

var _ layer.Base = &IncomingStack{}
