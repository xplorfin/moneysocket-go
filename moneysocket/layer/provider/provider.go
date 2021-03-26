package provider

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/provider"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

// ProviderLayer handles app waiting
// TODO this needs to be fully implemented
type ProviderLayer struct {
	layer.BaseLayer
	handlerProvideInfoRequest func(seed beacon.SharedSeed) account.AccountDb
	WaitingForApp             compat.WaitingForApp
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (o *ProviderLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(o.AnnounceNexus)
	belowLayer.SetOnRevoke(o.RevokeNexus)
}

// ProviderFinishedCb is the callback for the provider finished callback
func (o *ProviderLayer) ProviderFinishedCb(providerNexus nexus.Nexus) {
	o.TrackNexusAnnounced(providerNexus)
	o.SendLayerEvent(providerNexus, message.NexusAnnounced)
	o.OnAnnounce(providerNexus)
}

// AnnounceNexus creates a new provider.ProviderNexus and registers it
// also registers the providerFinishedCb (cb = callback)
func (o *ProviderLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	providerNexus := provider.NewProviderNexus(belowNexus)
	o.TrackNexus(providerNexus, belowNexus)
	providerNexus.WaitForConsumer(o.ProviderFinishedCb)
}

func (o *ProviderLayer) RevokeNexus(belowNexus nexus.Nexus) {
	res, _ := o.NexusByBelow.Get(belowNexus.UUID())
	providerNexus, _ := o.Nexuses.Get(res)
	o.BaseLayer.RevokeNexus(providerNexus)
	ss := providerNexus.SharedSeed()
	delete(o.WaitingForApp, ss.ToString())
}

func (o *ProviderLayer) HandlerProvideInfoRequest(seed beacon.SharedSeed) account.AccountDb {
	return o.handlerProvideInfoRequest(seed)
}

func (o *ProviderLayer) SetHandlerProvideInfoRequest(hpir compat.HandleProviderInfoRequest) {
	o.handlerProvideInfoRequest = hpir
}

func (o *ProviderLayer) NexusWaitingForApp(ss beacon.SharedSeed, providerNexus nexus.Nexus) {
	o.WaitingForApp[ss.ToString()] = providerNexus
}

func (o *ProviderLayer) ProviderNowReadyFromApp() {
	for sharedSeed, _ := range o.WaitingForApp { // nolint
		providerNexus := o.WaitingForApp[sharedSeed]
		delete(o.WaitingForApp, sharedSeed)
		providerNexus.(*provider.ProviderNexus).ProviderNowReady()
	}
}

func NewProviderLayer() *ProviderLayer {
	return &ProviderLayer{
		BaseLayer:     layer.NewBaseLayer(),
		WaitingForApp: make(map[string]nexus.Nexus),
	}
}

var _ layer.Layer = &ProviderLayer{}
