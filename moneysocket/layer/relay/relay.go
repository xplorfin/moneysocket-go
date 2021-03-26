package relay

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

type Layer struct {
	layer.BaseLayer
	RendezvousLayer *rendezvous.IncomingRendezvousLayer
}

func NewRelayLayer(rendezvousLayer *rendezvous.IncomingRendezvousLayer) *Layer {
	return &Layer{
		BaseLayer:       layer.NewBaseLayer(),
		RendezvousLayer: rendezvousLayer,
	}
}

// AnnounceNexus registers the message handlers for the rendezvousNexus to RelayLayer
func (r *Layer) AnnounceNexus(rendezvousNexus nexusHelper.Nexus) {
	rendezvousNexus.SetOnMessage(r.OnMessage)
	rendezvousNexus.SetOnBinMessage(r.OnBinMessage)
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (r *Layer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(r.AnnounceNexus)
	belowLayer.SetOnRevoke(r.RevokeNexus)
}

func (r *Layer) RevokeNexus(rendezvousNexus nexusHelper.Nexus) {
	log.Println("revoked from below")
}

func (r *Layer) OnMessage(rendezvousNexus nexusHelper.Nexus, msg base.MoneysocketMessage) {
	peerNexus := r.RendezvousLayer.GetPeerNexus(rendezvousNexus)
	_ = (*peerNexus).Send(msg)
}

func (r *Layer) OnBinMessage(rendezvousNexus nexusHelper.Nexus, msg []byte) {
	peerNexus := r.RendezvousLayer.GetPeerNexus(rendezvousNexus)
	_ = (*peerNexus).SendBin(msg)
}

var _ layer.Layer = &Layer{}
