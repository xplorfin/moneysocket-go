package terminus

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/local"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/provider"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	websocket2 "github.com/xplorfin/moneysocket-go/moneysocket/nexus/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/stack"
)

type OnStackEvent func(layerName string, nexus nexus.Nexus, status string)

type Stack struct {
	Config       *config.Config
	onAnnounce   layer.OnAnnounceFn
	onRevoke     layer.OnRevokeFn
	onStackEvent OnStackEvent

	localLayer      *local.OutgoingLocalLayer
	websocketLayer  *websocket.OutgoingWebsocketLayer
	rendezvousLayer *rendezvous.OutgoingRendezvousLayer
	providerLayer   *provider.Layer
	terminusLayer   *Layer
	incomingStack   *stack.IncomingStack

	// TODO add event listeners
}

func NewTerminusStack(config *config.Config) *Stack {
	s := Stack{Config: config}
	s.localLayer = s.SetupLocalLayer()
	s.websocketLayer = s.SetupWebsocketLayer()
	s.rendezvousLayer = s.SetupRendezvousLayer(s.websocketLayer, s.localLayer)
	s.providerLayer = s.SetupProviderLayer(s.rendezvousLayer)
	s.terminusLayer = s.SetupTerminusLayer(s.providerLayer)
	s.incomingStack = s.SetupIncomingStack(config, s.localLayer)

	return &s
}

func (s *Stack) SetupLocalLayer() *local.OutgoingLocalLayer {
	l := local.NewOutgoingLocalLayer()
	l.RegisterLayerEvent(s.SendStackEvent, message.OutgoingLocal)
	return &l
}

func (s *Stack) AnnounceNexus(terminusNexus Nexus) {
	if s.onAnnounce != nil {
		s.onAnnounce(&terminusNexus)
	}
}

func (s *Stack) SetupTerminusLayer(belowLayer layer.Layer) *Layer {
	l := NewTerminusLayer()
	l.RegisterAboveLayer(belowLayer)
	l.RegisterLayerEvent(s.onStackEvent, message.Terminus)
	// l.handleinvoicerequest = self.terminus_handle_invoice_request
	// l.handlepayrequest = self.terminus_handle_pay_request
	// l.handleproviderinforequest = self.terminus_handle_provider_info_request
	return l
}

func (s *Stack) SetupIncomingStack(config *config.Config, localLayer *local.OutgoingLocalLayer) *stack.IncomingStack {
	incomingStack := stack.NewIncomingStack(config, localLayer)
	return incomingStack
}

func (s *Stack) SetupWebsocketLayer() *websocket.OutgoingWebsocketLayer {
	l := websocket.NewOutgoingWebsocketLayer()
	l.RegisterLayerEvent(s.SendStackEvent, message.OutgoingWebsocket)
	return l
}

func (s *Stack) SetupProviderLayer(belowLayer layer.Layer) *provider.Layer {
	l := provider.NewProviderLayer()
	l.RegisterAboveLayer(belowLayer)
	l.RegisterLayerEvent(s.SendStackEvent, message.Provider)
	// l.handleproviderinforequest = handleproviderinforequest
	return l
}

func (s *Stack) SetupRendezvousLayer(belowLayer1 layer.Layer, belowLayer2 layer.Layer) *rendezvous.OutgoingRendezvousLayer {
	l := rendezvous.NewOutgoingRendezvousLayer()
	l.RegisterAboveLayer(belowLayer1)
	l.RegisterAboveLayer(belowLayer2)
	l.RegisterLayerEvent(s.SendStackEvent, message.OutgoingRendezvous)
	return l
}

func (s *Stack) SendStackEvent(layerName string, nexus nexus.Nexus, status string) {
	if s.onStackEvent != nil {
		s.onStackEvent(layerName, nexus, status)
	}
}

// get ws listen locations from the config
func (s *Stack) GetListenLocation() []location.Location {
	return s.incomingStack.GetListenLocations()
}

func (s *Stack) Connect(location location.WebsocketLocation, sharedSeed *beacon.SharedSeed) (*websocket2.OutgoingSocket, error) {
	return s.websocketLayer.Connect(location, sharedSeed)
}

func (s *Stack) LocalConnect(sharedSeed beacon.SharedSeed) {
	s.localLayer.Connect(sharedSeed)
}

func (s *Stack) Listen() error {
	return s.incomingStack.Listen()
}
