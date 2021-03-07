package provider

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/provider"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

// TODO this needs to be fully implemented
type ProviderLayer struct {
	layer.BaseLayer
	handlerProvideInfoRequest func(seed beacon.SharedSeed) account.AccountDb
	requestReferenceUuid      string
	providerFinishedCb        func(nexus2 nexus.Nexus)
	WaitingForApp             compat.WaitingForApp
}

func (o *ProviderLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(o.AnnounceNexus)
	belowLayer.SetOnRevoke(o.RevokeNexus)
}

func (o *ProviderLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	providerNexus := provider.NewProviderNexus(belowNexus)
	o.TrackNexus(providerNexus, belowNexus)
	providerNexus.WaitForConsumer(o.providerFinishedCb)
}

func (o *ProviderLayer) HandlerProvideInfoRequest(seed beacon.SharedSeed) account.AccountDb {
	return o.handlerProvideInfoRequest(seed)
}

func (o *ProviderLayer) SetHandlerProvideInfoRequest(hpir compat.HandleProviderInfoRequest) {
	o.handlerProvideInfoRequest = hpir
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
