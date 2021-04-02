package relay

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// Layer is the relay layer.
type Layer struct {
	layer.BaseLayer
	// RendezvousLayer is the rendezvous.IncomingRendezvousLayer that hits a peer nexus
	RendezvousLayer *rendezvous.IncomingRendezvousLayer
}

// NewRelayLayer creates the Layer given the rendezvous layer.
func NewRelayLayer(rendezvousLayer *rendezvous.IncomingRendezvousLayer) *Layer {
	return &Layer{
		BaseLayer:       layer.NewBaseLayer(),
		RendezvousLayer: rendezvousLayer,
	}
}

// AnnounceNexus registers the message handlers for the rendezvousNexus to RelayLayer.
func (r *Layer) AnnounceNexus(rendezvousNexus nexusHelper.Nexus) {
	rendezvousNexus.SetOnMessage(r.OnMessage)
	rendezvousNexus.SetOnBinMessage(r.OnBinMessage)
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer.
func (r *Layer) RegisterAboveLayer(belowLayer layer.Base) {
	belowLayer.SetOnAnnounce(r.AnnounceNexus)
	belowLayer.SetOnRevoke(r.RevokeNexus)
}

// RevokeNexus revokes a nexus from below. This method is simply used to conform to the layer.Base here.
func (r *Layer) RevokeNexus(rendezvousNexus nexusHelper.Nexus) {
	log.Println("revoked from below")
}

// OnMessage sends a non-binary message to the peer nexus.
func (r *Layer) OnMessage(rendezvousNexus nexusHelper.Nexus, msg base.MoneysocketMessage) {
	peerNexus := r.RendezvousLayer.GetPeerNexus(rendezvousNexus)
	_ = (*peerNexus).Send(msg)
}

// OnBinMessage takes a binary message and sends it to the peer nexus.
func (r *Layer) OnBinMessage(rendezvousNexus nexusHelper.Nexus, msg []byte) {
	peerNexus := r.RendezvousLayer.GetPeerNexus(rendezvousNexus)
	_ = (*peerNexus).SendBin(msg)
}

var _ layer.Base = &Layer{}
