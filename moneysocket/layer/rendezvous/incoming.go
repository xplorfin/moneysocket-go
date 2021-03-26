package rendezvous

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/rendezvous"
)

// IncomingRendezvousLayer is responsible for managing rendezvous' between peernexuses
type IncomingRendezvousLayer struct {
	*layer.BaseLayer
	// directory is used for peering nexuses
	directory *rendezvous.Directory
}

// NewIncomingRendezvousLayer creates an IncomingRendezvousLayer
func NewIncomingRendezvousLayer() *IncomingRendezvousLayer {
	baseLayer := layer.NewBaseLayer()
	il := IncomingRendezvousLayer{
		BaseLayer: &baseLayer,
		directory: rendezvous.NewRendezvousDirectory(),
	}
	il.SetOnAnnounce(il.AnnounceNexus)
	il.SetOnRevoke(il.RevokeNexus)
	return &il
}

// AnnounceNexus creates a new rendezvous.IncomingRendezvousNexus and registers it
func (o *IncomingRendezvousLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	rendezvousNexus := rendezvous.NewIncomingRendezvousNexus(belowNexus, o, o.directory)

	o.TrackNexus(rendezvousNexus, belowNexus)
	rendezvousNexus.WaitForRendezvous(o.RendezvousFinishedCb)
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (o *IncomingRendezvousLayer) RegisterAboveLayer(belowLayer layer.LayerBase) {
	belowLayer.SetOnAnnounce(o.AnnounceNexus)
	belowLayer.SetOnRevoke(o.RevokeNexus)
}

// RendezvousFinishedCb is the callback for after a rendezvous is finished
func (o *IncomingRendezvousLayer) RendezvousFinishedCb(rendezvousNexus nexus.Nexus) {
	o.TrackNexusAnnounced(rendezvousNexus)
	o.SendLayerEvent(rendezvousNexus, message.NexusAnnounced)
	if o.OnAnnounce != nil {
		o.OnAnnounce(rendezvousNexus)
	}
}

func (o *IncomingRendezvousLayer) ToString() string {
	return o.directory.ToString()
}

// RevokeNexus removes the nexus from directories/layers
func (o *IncomingRendezvousLayer) RevokeNexus(belowNexus nexus.Nexus) {
	belowUUID, _ := o.NexusByBelow.Get(belowNexus.UUID())
	rendezvousNexus, _ := o.Nexuses.Get(belowUUID)
	peerRendezvousNexus := o.directory.GetPeerNexus(rendezvousNexus.UUID())
	o.BaseLayer.RevokeNexus(belowNexus)
	if peerRendezvousNexus != nil {
		o.directory.RemoveNexus(*peerRendezvousNexus)
		irNexus := (*peerRendezvousNexus).(*rendezvous.IncomingRendezvousNexus)
		irNexus.EndRendezvous()
	}
}

// GetPeerNexus is the gets the peered nexus from the directory for a nexus (using uuid)
func (o *IncomingRendezvousLayer) GetPeerNexus(rendezvousNexus nexus.Nexus) *nexus.Nexus {
	return o.directory.GetPeerNexus(rendezvousNexus.UUID())
}

var _ layer.LayerBase = &IncomingRendezvousLayer{}
