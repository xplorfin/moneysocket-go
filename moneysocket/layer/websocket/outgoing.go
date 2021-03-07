package websocket

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	nws "github.com/xplorfin/moneysocket-go/moneysocket/nexus/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/ws/ws_client"
)

// TODO this needs to be fully implemented
type OutgoingWebsocketLayer struct {
	layer.BaseLayer
	NexusBySharedSeed      layer.NexusStringMap
	OutgoingSocketProtocol *nws.OutgoingSocket
}

func NewOutgoingWebsocketLayer() *OutgoingWebsocketLayer {
	outgoingSocket := nws.NewOutgoingSocket()
	os := OutgoingWebsocketLayer{
		BaseLayer:              layer.NewBaseLayer(),
		NexusBySharedSeed:      layer.NexusStringMap{},
		OutgoingSocketProtocol: &outgoingSocket,
	}
	os.OutgoingSocketProtocol.FactoryMsProtocolLayer = &os
	return &os
}

func (o *OutgoingWebsocketLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	websocketNexus := nws.NewWebsocketNexus(belowNexus, o)
	o.TrackNexus(&websocketNexus, belowNexus)
	o.TrackNexusAnnounced(&websocketNexus)

	sharedSeed := websocketNexus.SharedSeed()
	o.NexusBySharedSeed.Store(sharedSeed.ToString(), &websocketNexus)
	o.SendLayerEvent(&websocketNexus, message.NexusAnnounced)
	if o.OnAnnounce != nil {
		o.OnAnnounce(&websocketNexus)
	}
}

func (o *OutgoingWebsocketLayer) Connect(location location.WebsocketLocation, seed *beacon.SharedSeed) (*nws.OutgoingSocket, error) {
	o.OutgoingSocketProtocol.FactoryMsProtocolLayer = o
	o.OutgoingSocketProtocol.OutgoingSharedSeed = seed
	// we do this in a func so a traceback leads back here
	go func() {
		ws_client.NewWsClient(o.OutgoingSocketProtocol, location.ToString())
	}()
	return o.OutgoingSocketProtocol, nil
}

func (o *OutgoingWebsocketLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	panic("implement me")
}

func (o *OutgoingWebsocketLayer) InitiateCloseAll() {
	o.Nexuses.Range(func(key uuid.UUID, nexus nexus.Nexus) bool {
		nexus.InitiateClose()
		return true
	})
}

var _ layer.Layer = &OutgoingWebsocketLayer{}
