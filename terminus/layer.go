package terminus

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// TODO this needs to be fully implemented
type TerminusLayer struct {
	layer.BaseLayer
	NexusesBySharedSeed       map[string][]string
	handleInvoiceRequest      compat.HandleInvoiceRequest
	handlePayRequest          compat.HandlePayRequest
	handleProviderInfoRequest compat.HandleProviderInfoRequest
}

func (o *TerminusLayer) SetupTerminusNexus(belowNexus nexus.Nexus) *TerminusNexus {
	terminusNexus := NewTerminusNexus(belowNexus, o)
	terminusNexus.handleProviderInfoRequest = o.handleProviderInfoRequest
	terminusNexus.handleInvoiceRequest = o.handleInvoiceRequest
	terminusNexus.handlePayRequest = o.handlePayRequest
	return &terminusNexus
}

func (o *TerminusLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	terminusNexus := o.SetupTerminusNexus(belowNexus)
	o.TrackNexus(terminusNexus, belowNexus)
	o.SendLayerEvent(terminusNexus, message.NexusAnnounced)
	if o.OnAnnounce != nil {
		o.OnAnnounce(terminusNexus)
	}

	ss := terminusNexus.SharedSeed()
	if _, ok := o.NexusesBySharedSeed[ss.ToString()]; !ok {
		o.NexusesBySharedSeed[ss.ToString()] = []string{}
	}
	o.NexusesBySharedSeed[ss.ToString()] = append(o.NexusesBySharedSeed[ss.ToString()], terminusNexus.Uuid().String())
}

func (o *TerminusLayer) RevokeNexus(belowNexus nexus.Nexus) {
	nexusUuid, _ := o.NexusByBelow.Get(belowNexus.Uuid())
	terminusNexus, _ := o.Nexuses.Get(nexusUuid)
	o.BaseLayer.RevokeNexus(terminusNexus)
	ss := terminusNexus.SharedSeed()
	delete(o.NexusesBySharedSeed, ss.ToString())
}

func (o *TerminusLayer) NotifyPreImage(sharedSeeds []beacon.SharedSeed, preimage string) {
	for _, ss := range sharedSeeds {
		if _, ok := o.NexusesBySharedSeed[ss.ToString()]; !ok {
			continue
		}
		for _, nexusUuid := range o.NexusesBySharedSeed[ss.ToString()] {
			nxId, _ := uuid.FromString(nexusUuid)
			nx, _ := o.Nexuses.Get(nxId)
			terminusNexus := nx.(*TerminusNexus)
			terminusNexus.NotifyPreimage(preimage, uuid.NewV4().String())
			terminusNexus.NotifyProviderInfo(ss)
		}
	}
}

func (o *TerminusLayer) HandlePayRequest(ss beacon.SharedSeed, bolt11 string) {
	panic("method not yet implemented")
}

func (o *TerminusLayer) HandleInvoiceRequest(ss beacon.SharedSeed, msats int) {
	panic("method not yet implemented")
}

func (o *TerminusLayer) HandleProviderInfoRequest(ss beacon.SharedSeed) {
	panic("method not yet implemented")
}

func (o *TerminusLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(o.OnAnnounce)
	belowLayer.SetOnRevoke(o.OnRevoke)
}

func NewTerminusLayer() *TerminusLayer {
	return &TerminusLayer{
		BaseLayer: layer.NewBaseLayer(),
	}
}

var _ layer.Layer = &TerminusLayer{}
