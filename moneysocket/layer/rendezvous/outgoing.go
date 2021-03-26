package rendezvous

import (
	"encoding/hex"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/rendezvous"
)

// TODO this needs to be fully implemented
type OutgoingRendezvousLayer struct {
	layer.BaseLayer
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (o *OutgoingRendezvousLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(o.AnnounceNexus)
	belowLayer.SetOnRevoke(o.RevokeNexus)
}

// AnnounceNexus creates a new rendezvous.OutgoingRendezvousNexus and registers it
// a rendezvous is started and if completed OutgoingRendezvousLayer.RendezvousFinishedCb is called
func (o *OutgoingRendezvousLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	rendezvousNexus := rendezvous.NewOutgoingRendezvousNexus(belowNexus, o)
	o.TrackNexus(rendezvousNexus, belowNexus)
	sharedSeed := belowNexus.SharedSeed()
	rid := hex.EncodeToString(sharedSeed.DeriveRendezvousID())
	rendezvousNexus.StartRendezvous(rid, o.RendezvousFinishedCb)
}

func (o *OutgoingRendezvousLayer) RendezvousFinishedCb(rendzvousNexus nexus.Nexus) {
	o.TrackNexusAnnounced(rendzvousNexus)
	o.SendLayerEvent(rendzvousNexus, message.NexusAnnounced)
	if o.OnAnnounce != nil {
		o.OnAnnounce(rendzvousNexus)
	}
}

func NewOutgoingRendezvousLayer() *OutgoingRendezvousLayer {
	return &OutgoingRendezvousLayer{
		BaseLayer: layer.NewBaseLayer(),
	}
}

var _ layer.Layer = &OutgoingRendezvousLayer{}
