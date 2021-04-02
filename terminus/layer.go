package terminus

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// Layer is the terminus layer TODO this needs to be fully implemented.
type Layer struct {
	layer.BaseLayer
	NexusesBySharedSeed       map[string][]string
	handleInvoiceRequest      compat.HandleInvoiceRequest
	handlePayRequest          compat.HandlePayRequest
	handleProviderInfoRequest compat.HandleProviderInfoRequest
}

// SetupTerminusNexus sets up the terminus nexus.
func (o *Layer) SetupTerminusNexus(belowNexus nexus.Nexus) *Nexus {
	terminusNexus := NewTerminusNexus(belowNexus, o)
	terminusNexus.handleProviderInfoRequest = o.handleProviderInfoRequest
	terminusNexus.handleInvoiceRequest = o.handleInvoiceRequest
	terminusNexus.handlePayRequest = o.handlePayRequest
	return &terminusNexus
}

// AnnounceNexus creates a new TerminusNexus and registers it.
func (o *Layer) AnnounceNexus(belowNexus nexus.Nexus) {
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
	o.NexusesBySharedSeed[ss.ToString()] = append(o.NexusesBySharedSeed[ss.ToString()], terminusNexus.UUID().String())
}

// RevokeNexus removes a nexus.
func (o *Layer) RevokeNexus(belowNexus nexus.Nexus) {
	nexusUUID, _ := o.NexusByBelow.Get(belowNexus.UUID())
	terminusNexus, _ := o.Nexuses.Get(nexusUUID)
	o.BaseLayer.RevokeNexus(terminusNexus)
	ss := terminusNexus.SharedSeed()
	delete(o.NexusesBySharedSeed, ss.ToString())
}

// NotifyPreImage notifies a preimage.
func (o *Layer) NotifyPreImage(sharedSeeds []beacon.SharedSeed, preimage string) {
	for _, ss := range sharedSeeds {
		if _, ok := o.NexusesBySharedSeed[ss.ToString()]; !ok {
			continue
		}
		for _, nexusUUID := range o.NexusesBySharedSeed[ss.ToString()] {
			nxID, _ := uuid.FromString(nexusUUID)
			nx, _ := o.Nexuses.Get(nxID)
			terminusNexus := nx.(*Nexus)
			terminusNexus.NotifyPreimage(preimage, uuid.NewV4().String())
			terminusNexus.NotifyProviderInfo(ss)
		}
	}
}

// HandlePayRequest handles a payment request.
func (o *Layer) HandlePayRequest(ss beacon.SharedSeed, bolt11 string) {
	panic("method not yet implemented")
}

// HandleInvoiceRequest handles an invoice request.
func (o *Layer) HandleInvoiceRequest(ss beacon.SharedSeed, msats int) {
	panic("method not yet implemented")
}

// HandleProviderInfoRequest processes a provider info request.
func (o *Layer) HandleProviderInfoRequest(ss beacon.SharedSeed) {
	panic("method not yet implemented")
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer.
func (o *Layer) RegisterAboveLayer(belowLayer layer.Base) {
	belowLayer.SetOnAnnounce(o.OnAnnounce)
	belowLayer.SetOnRevoke(o.OnRevoke)
}

// NewTerminusLayer creates a new terminus layer.
func NewTerminusLayer() *Layer {
	return &Layer{
		BaseLayer: layer.NewBaseLayer(),
	}
}

var _ layer.Base = &Layer{}
