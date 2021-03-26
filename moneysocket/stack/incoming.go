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

// TODO
type IncomingStack struct {
	layer.BaseLayer
	config          *config.Config
	localLayer      *local.IncomingLocalLayer
	websocketLayer  *websocket.IncomingWebsocketLayer
	rendezvousLayer *rendezvous.IncomingRendezvousLayer
	relayLayer      *relay.Layer
}

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

func (i *IncomingStack) SetupRelayLayer(rendezvousLayer *rendezvous.IncomingRendezvousLayer) {
	i.relayLayer = relay.NewRelayLayer(rendezvousLayer)
	i.relayLayer.RegisterAboveLayer(rendezvousLayer)
	i.relayLayer.RegisterLayerEvent(i.SendStackEvent, message.Relay)
}

func (i *IncomingStack) SetupRendezvousLayer(belowLayer1 layer.Layer, belowLayer2 layer.Layer) {
	i.rendezvousLayer = rendezvous.NewIncomingRendezvousLayer()
	i.rendezvousLayer.RegisterAboveLayer(belowLayer1)
	i.rendezvousLayer.RegisterAboveLayer(belowLayer2)
	i.rendezvousLayer.RegisterLayerEvent(i.SendStackEvent, message.IncomingRendezvous)
}

func (i *IncomingStack) SetupWebsocketLayer() {
	i.websocketLayer = websocket.NewIncomingWebsocketLayer(i.config)
	i.websocketLayer.RegisterLayerEvent(i.SendStackEvent, message.IncomingWebsocket)
}

func (i *IncomingStack) SetupLocalLayer(outgoingLocalLayer *local.OutgoingLocalLayer) {
	i.localLayer = local.NewIncomingLocalLayer()
	i.localLayer.RegisterLayerEvent(i.SendStackEvent, message.IncomingLocal)
	outgoingLocalLayer.SetIncomingLayer(i.localLayer)
}

func (i *IncomingStack) SendStackEvent(layerName string, nexus nexusHelper.Nexus, event string) {
	// do nothing
}

// RegisterAboveLayer does nothing. It is here to satisfy the interface since IncomingStack must
// act as a layer
func (i *IncomingStack) RegisterAboveLayer(belowLayer layer.Layer) {
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

// return websocket location object from config
func (i *IncomingStack) GetListenLocations() []location.Location {
	return []location.Location{location.NewWebsocketLocationPort(i.config.GetExternalHost(), i.config.GetUseTLS(), i.config.GetExternalPort())}
}

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

var _ layer.Layer = &IncomingStack{}
