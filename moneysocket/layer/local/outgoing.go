package local

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/local"
)

type OutgoingLocalLayer struct {
	layer.BaseLayer
	IncomingLocalLayer   layer.LayerBase
	OutgoingBySharedSeed map[string]nexus.Nexus
	IncomingBySharedSeed map[string]nexus.Nexus
}

func NewOutgoingLocalLayer() OutgoingLocalLayer {
	return OutgoingLocalLayer{
		BaseLayer:            layer.NewBaseLayer(),
		OutgoingBySharedSeed: make(map[string]nexus.Nexus),
		IncomingBySharedSeed: make(map[string]nexus.Nexus),
	}
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (o *OutgoingLocalLayer) RegisterAboveLayer(belowLayer layer.LayerBase) {
	belowLayer.SetOnAnnounce(o.AnnounceNexus)
	belowLayer.SetOnRevoke(o.RevokeNexus)
}

// AnnounceNexus creates a new local.LocalNexus and registers it
func (o *OutgoingLocalLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	localNexus := local.NewLocalNexus(belowNexus, o)

	o.TrackNexus(localNexus, belowNexus)
	o.TrackNexusAnnounced(localNexus)
	o.SendLayerEvent(localNexus, message.NexusAnnounced)
	if o.OnAnnounce != nil {
		o.OnAnnounce(localNexus)
	}
}

func (o *OutgoingLocalLayer) SetIncomingLayer(incomingLayer layer.LayerBase) {
	o.IncomingLocalLayer = incomingLayer
}

func (o *OutgoingLocalLayer) Connect(sharedSeed beacon.SharedSeed) {
	joinedNexus := local.NewJoinedLocalNexus()
	outgoingNexus := local.NewOutgoingLocalNexus(joinedNexus, o, sharedSeed)
	o.OutgoingBySharedSeed[sharedSeed.ToString()] = outgoingNexus
	incomingNexus := local.NewIncomingLocalNexus(joinedNexus, o.IncomingLocalLayer)
	// add incoming nexus on message
	o.IncomingBySharedSeed[sharedSeed.ToString()] = incomingNexus

	o.IncomingLocalLayer.AnnounceNexus(incomingNexus)

	o.AnnounceNexus(outgoingNexus)
}

func (o *OutgoingLocalLayer) Disconnect(sharedSeed beacon.SharedSeed) {
	outgoingNexus := o.OutgoingBySharedSeed[sharedSeed.ToString()]
	incomingNexus := o.IncomingBySharedSeed[sharedSeed.ToString()]
	delete(o.OutgoingBySharedSeed, sharedSeed.ToString())
	delete(o.IncomingBySharedSeed, sharedSeed.ToString())
	o.IncomingLocalLayer.RevokeNexus(incomingNexus)
	o.RevokeNexus(outgoingNexus)
}

// statically assert outgoing layer matches layer interface
var _ layer.LayerBase = &OutgoingLocalLayer{}
