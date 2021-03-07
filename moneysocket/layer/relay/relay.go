package relay

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

type RelayLayer struct {
	layer.BaseLayer
	RendezvousLayer *rendezvous.IncomingRendezvousLayer
}

func NewRelayLayer(rendezvousLayer *rendezvous.IncomingRendezvousLayer) *RelayLayer {
	return &RelayLayer{
		BaseLayer:       layer.NewBaseLayer(),
		RendezvousLayer: rendezvousLayer,
	}
}

func (r *RelayLayer) AnnounceNexus(rendezvousNexus nexusHelper.Nexus) {
	rendezvousNexus.SetOnMessage(r.OnMessage)
	rendezvousNexus.SetOnBinMessage(r.OnBinMessage)
}

func (r *RelayLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(r.AnnounceNexus)
	belowLayer.SetOnRevoke(r.RevokeNexus)
}

func (r *RelayLayer) RevokeNexus(rendezvousNexus nexusHelper.Nexus) {
	log.Println("revoked from below")
}

func (r *RelayLayer) OnMessage(rendezvousNexus nexusHelper.Nexus, msg base.MoneysocketMessage) {
	peerNexus := r.RendezvousLayer.GetPeerNexus(rendezvousNexus)
	_ = (*peerNexus).Send(msg)
}

func (r *RelayLayer) OnBinMessage(rendezvousNexus nexusHelper.Nexus, msg []byte) {
	peerNexus := r.RendezvousLayer.GetPeerNexus(rendezvousNexus)
	_ = (*peerNexus).SendBin(msg)
}

var _ layer.Layer = &RelayLayer{}
