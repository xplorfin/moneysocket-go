package local

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/local"
)

type IncomingLocalLayer struct {
	layer.BaseLayer
}

func NewIncomingLocalLayer() *IncomingLocalLayer {
	return &IncomingLocalLayer{
		BaseLayer: layer.NewBaseLayer(),
	}
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (i *IncomingLocalLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(i.AnnounceNexus)
	belowLayer.SetOnRevoke(i.RevokeNexus)
}

// AnnounceNexus creates a new LocalNexus and registers it
func (i *IncomingLocalLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	localNexus := local.NewLocalNexus(belowNexus, i)
	// register above nexus
	i.TrackNexus(localNexus, belowNexus)
	i.TrackNexusAnnounced(belowNexus)
	i.SendLayerEvent(localNexus, message.NexusAnnounced)
	if i.OnAnnounce != nil {
		i.OnAnnounce(localNexus)
	}
}
