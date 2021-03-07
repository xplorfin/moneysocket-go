package provider

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/provider"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

type ProviderTransactLayer struct {
	layer.BaseLayer
	handleInvoiceRequest      compat.HandleInvoiceRequest
	handlePayRequest          compat.HandlePayRequest
	handleProviderInfoRequest compat.HandleProviderInfoRequest
	NexusBySharedSeed         compat.NexusBySharedSeed
}

func (p *ProviderTransactLayer) HandleProviderInfoRequest(seed beacon.SharedSeed) account.AccountDb {
	return p.handleProviderInfoRequest(seed)
}

func (p *ProviderTransactLayer) HandlePayRequest(nexus nexus.Nexus, bolt11 string, requestUuid string) {
	p.handlePayRequest(nexus, bolt11, requestUuid)
}

func (p *ProviderTransactLayer) HandleInvoiceRequest(nexus nexus.Nexus, msats int64, requestUuid string) {
	p.handleInvoiceRequest(nexus, msats, requestUuid)
}

func (p *ProviderTransactLayer) SetHandleInvoiceRequest(request compat.HandleInvoiceRequest) {
	p.handleInvoiceRequest = request
}

func (p *ProviderTransactLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(p.AnnounceNexus)
	belowLayer.SetOnRevoke(p.OnRevoke)
}

func NewProviderTransactLayer() *ProviderTransactLayer {
	return &ProviderTransactLayer{
		BaseLayer:         layer.NewBaseLayer(),
		NexusBySharedSeed: make(compat.NexusBySharedSeed),
	}
}

func (p *ProviderTransactLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	// setup the transaction nexus
	providerTransactNexus := p.setupTransactNexus(belowNexus)
	p.TrackNexus(providerTransactNexus, belowNexus)
	p.TrackNexusAnnounced(providerTransactNexus)

	if p.OnAnnounce != nil {
		p.OnAnnounce(providerTransactNexus)
	}

	sharedSeed := providerTransactNexus.SharedSeed()
	if _, ok := p.NexusBySharedSeed[sharedSeed.ToString()]; !ok {
		p.NexusBySharedSeed[sharedSeed.ToString()] = []nexus.Nexus{}
	}
	p.NexusBySharedSeed[sharedSeed.ToString()] = append(p.NexusBySharedSeed[sharedSeed.ToString()], providerTransactNexus)
}

// create a transact nexus
func (p *ProviderTransactLayer) setupTransactNexus(belowNexus nexus.Nexus) *provider.ProviderTransactNexus {
	providerTransactNexus := provider.NewProviderTransactNexus(belowNexus, p)
	providerTransactNexus.HandleInvoiceRequest = p.handleInvoiceRequest
	providerTransactNexus.HandlePayRequest = p.handlePayRequest
	return providerTransactNexus
}

func (p *ProviderTransactLayer) RevokeNexus(belowNexus nexus.Nexus) {
	belowUuid, _ := p.NexusByBelow.Get(belowNexus.Uuid())
	providerTransactNexus, _ := p.Nexuses.Get(belowUuid)
	p.BaseLayer.RevokeNexus(providerTransactNexus)
	sharedSeed := providerTransactNexus.SharedSeed()
	var nexusIndex int
	for i, nexusSeed := range p.NexusBySharedSeed[sharedSeed.ToString()] {
		if nexusSeed.Uuid() == providerTransactNexus.Uuid() {
			nexusIndex = i
		}
	}
	p.NexusBySharedSeed[sharedSeed.ToString()] = append(p.NexusBySharedSeed[sharedSeed.ToString()][:nexusIndex], p.NexusBySharedSeed[sharedSeed.ToString()][nexusIndex:]...)
}

func (p *ProviderTransactLayer) FulfilRequestInvoice(nexusUuid, bolt11, requestReferenceUuid string) error {
	nexusId, _ := uuid.FromBytes([]byte(nexusUuid))
	if nx, ok := p.Nexuses.Get(nexusId); ok {
		providerTxNexus := nx.(*provider.ProviderTransactNexus)
		return providerTxNexus.NotifyInvoice(bolt11, requestReferenceUuid)
	}
	return nil
}

func (p *ProviderTransactLayer) NotifyPreImage(sharedSeeds []beacon.SharedSeed, preimage, requestReferenceUuid string) {
	for _, sharedSeed := range sharedSeeds {
		if _, ok := p.NexusBySharedSeed[sharedSeed.ToString()]; ok {
			continue
		}
		for _, nx := range p.NexusBySharedSeed[sharedSeed.ToString()] {
			ptn := nx.(*provider.ProviderTransactNexus)
			_ = ptn.NotifyPreimage(preimage, requestReferenceUuid)
		}
	}
}

func (p *ProviderTransactLayer) NotifyProviderInfo(sharedSeeds []beacon.SharedSeed, preimage, requestReferenceUuid string) {
	for _, sharedSeed := range sharedSeeds {
		if _, ok := p.NexusBySharedSeed[sharedSeed.ToString()]; ok {
			continue
		}
		for _, nx := range p.NexusBySharedSeed[sharedSeed.ToString()] {
			ptn := nx.(*provider.ProviderTransactNexus)
			_ = ptn.NotifyProviderInfo(sharedSeed)
		}
	}
}

var _ layer.Layer = &ProviderTransactLayer{}

// use an interface to call methods in the nexus
var _ compat.ProviderTransactLayerInterface = &ProviderTransactLayer{}
