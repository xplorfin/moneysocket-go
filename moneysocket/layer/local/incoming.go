package local

import (
	"fmt"

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

func (i *IncomingLocalLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(i.AnnounceNexus)
	belowLayer.SetOnRevoke(i.RevokeNexus)
}

func (i *IncomingLocalLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	localNexus := local.NewLocalNexus(belowNexus, i)
	// register above nexus
	belowNexus.SetOnMessage(localNexus.OnMessage)
	belowNexus.SetOnBinMessage(localNexus.OnBinMessage)
	i.TrackNexus(&localNexus, belowNexus)
	i.TrackNexusAnnounced(belowNexus)
	i.SendLayerEvent(&localNexus, message.NexusAnnounced)
	if i.OnAnnounce != nil {
		i.OnAnnounce(&localNexus)
	} else {
		fmt.Println("hi")
	}
}
