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

type TransactLayer struct {
	layer.BaseLayer
	handleInvoiceRequest      compat.HandleInvoiceRequest
	handlePayRequest          compat.HandlePayRequest
	handleProviderInfoRequest compat.HandleProviderInfoRequest
	NexusBySharedSeed         compat.NexusBySharedSeed
}

func (p *TransactLayer) HandleProviderInfoRequest(seed beacon.SharedSeed) account.Db {
	return p.handleProviderInfoRequest(seed)
}

func (p *TransactLayer) HandlePayRequest(nexus nexus.Nexus, bolt11, requestUUID string) {
	p.handlePayRequest(nexus, bolt11, requestUUID)
}

func (p *TransactLayer) HandleInvoiceRequest(nexus nexus.Nexus, msats int64, requestUUID string) {
	p.handleInvoiceRequest(nexus, msats, requestUUID)
}

func (p *TransactLayer) SetHandleInvoiceRequest(request compat.HandleInvoiceRequest) {
	p.handleInvoiceRequest = request
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (p *TransactLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(p.AnnounceNexus)
	belowLayer.SetOnRevoke(p.OnRevoke)
}

func NewProviderTransactLayer() *TransactLayer {
	return &TransactLayer{
		BaseLayer:         layer.NewBaseLayer(),
		NexusBySharedSeed: make(compat.NexusBySharedSeed),
	}
}

func (p *TransactLayer) AnnounceNexus(belowNexus nexus.Nexus) {
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
func (p *TransactLayer) setupTransactNexus(belowNexus nexus.Nexus) *provider.TransactNexus {
	providerTransactNexus := provider.NewProviderTransactNexus(belowNexus, p)
	providerTransactNexus.HandleInvoiceRequest = p.handleInvoiceRequest
	providerTransactNexus.HandlePayRequest = p.handlePayRequest
	return providerTransactNexus
}

func (p *TransactLayer) RevokeNexus(belowNexus nexus.Nexus) {
	belowUUID, _ := p.NexusByBelow.Get(belowNexus.UUID())
	providerTransactNexus, _ := p.Nexuses.Get(belowUUID)
	p.BaseLayer.RevokeNexus(providerTransactNexus)
	sharedSeed := providerTransactNexus.SharedSeed()
	var nexusIndex int
	for i, nexusSeed := range p.NexusBySharedSeed[sharedSeed.ToString()] {
		if nexusSeed.UUID() == providerTransactNexus.UUID() {
			nexusIndex = i
		}
	}
	p.NexusBySharedSeed[sharedSeed.ToString()] = append(p.NexusBySharedSeed[sharedSeed.ToString()][:nexusIndex], p.NexusBySharedSeed[sharedSeed.ToString()][nexusIndex:]...)
}

func (p *TransactLayer) FulfilRequestInvoice(nexusUUID, bolt11, requestReferenceUUID string) error {
	nexusID, _ := uuid.FromBytes([]byte(nexusUUID))
	if nx, ok := p.Nexuses.Get(nexusID); ok {
		providerTxNexus := nx.(*provider.TransactNexus)
		return providerTxNexus.NotifyInvoice(bolt11, requestReferenceUUID)
	}
	return nil
}

func (p *TransactLayer) NotifyPreImage(sharedSeeds []beacon.SharedSeed, preimage, requestReferenceUUID string) {
	for _, sharedSeed := range sharedSeeds {
		if _, ok := p.NexusBySharedSeed[sharedSeed.ToString()]; ok {
			continue
		}
		for _, nx := range p.NexusBySharedSeed[sharedSeed.ToString()] {
			ptn := nx.(*provider.TransactNexus)
			_ = ptn.NotifyPreimage(preimage, requestReferenceUUID)
		}
	}
}

func (p *TransactLayer) NotifyProviderInfo(sharedSeeds []beacon.SharedSeed) {
	for _, sharedSeed := range sharedSeeds {
		if _, ok := p.NexusBySharedSeed[sharedSeed.ToString()]; ok {
			continue
		}
		for _, nx := range p.NexusBySharedSeed[sharedSeed.ToString()] {
			ptn := nx.(*provider.TransactNexus)
			_ = ptn.NotifyProviderInfo(sharedSeed)
		}
	}
}

var _ layer.Layer = &TransactLayer{}

// use an interface to call methods in the nexus
var _ compat.ProviderTransactLayerInterface = &TransactLayer{}
