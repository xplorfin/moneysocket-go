package rendezvous

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/rendezvous"
)

type IncomingRendezvousLayer struct {
	*layer.BaseLayer
	directory *rendezvous.RendezvousDirectory
}

func (o *IncomingRendezvousLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(o.AnnounceNexus)
	belowLayer.SetOnRevoke(o.RevokeNexus)
}

func (o *IncomingRendezvousLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	rendezvousNexus := rendezvous.NewIncomingRendezvousNexus(belowNexus, o, o.directory)

	o.TrackNexus(&rendezvousNexus, belowNexus)
	rendezvousNexus.WaitForRendezvous(o.RendezvousFinishedCb)
}

func (o *IncomingRendezvousLayer) RendezvousFinishedCb(rendezvousNexus nexus.Nexus) {
	o.TrackNexusAnnounced(rendezvousNexus)
	o.SendLayerEvent(rendezvousNexus, message.NexusAnnounced)
	if o.OnAnnounce != nil {
		o.OnAnnounce(rendezvousNexus)
	}
}

func (o *IncomingRendezvousLayer) RevokeNexus(belowNexus nexus.Nexus) {
	belowUuid, _ := o.NexusByBelow.Get(belowNexus.Uuid())
	rendezvousNexus, _ := o.Nexuses.Get(belowUuid)
	peerRendezvousNexus := o.directory.GetPeerNexus(rendezvousNexus.Uuid())
	o.BaseLayer.RevokeNexus(belowNexus)
	if peerRendezvousNexus != nil {
		o.directory.RemoveNexus(*peerRendezvousNexus)
		irNexus := (*peerRendezvousNexus).(*rendezvous.IncomingRendezvousNexus)
		irNexus.EndRendezvous()
	}
}

func (o *IncomingRendezvousLayer) GetPeerNexus(rendezvousNexus nexus.Nexus) nexus.Nexus {
	panic("method not yet implemented")
}

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

var _ layer.Layer = &IncomingRendezvousLayer{}
