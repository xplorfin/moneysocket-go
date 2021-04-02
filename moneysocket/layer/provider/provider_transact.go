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

// TransactLayer handles transactions.
type TransactLayer struct {
	layer.BaseLayer
	handleInvoiceRequest      compat.HandleInvoiceRequest
	handlePayRequest          compat.HandlePayRequest
	handleProviderInfoRequest compat.HandleProviderInfoRequest
	// NexusBySharedSeed tracks nexuses by shared seed
	NexusBySharedSeed compat.NexusBySharedSeed
}

// HandleProviderInfoRequest handles info requests.
func (p *TransactLayer) HandleProviderInfoRequest(seed beacon.SharedSeed) account.DB {
	return p.handleProviderInfoRequest(seed)
}

// HandlePayRequest handles payment requests.
func (p *TransactLayer) HandlePayRequest(nexus nexus.Nexus, bolt11, requestUUID string) {
	p.handlePayRequest(nexus, bolt11, requestUUID)
}

// HandleInvoiceRequest gets an invoice.
func (p *TransactLayer) HandleInvoiceRequest(nexus nexus.Nexus, msats int64, requestUUID string) {
	p.handleInvoiceRequest(nexus, msats, requestUUID)
}

// SetHandleInvoiceRequest sets an invoice request.
func (p *TransactLayer) SetHandleInvoiceRequest(request compat.HandleInvoiceRequest) {
	p.handleInvoiceRequest = request
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer.
func (p *TransactLayer) RegisterAboveLayer(belowLayer layer.Base) {
	belowLayer.SetOnAnnounce(p.AnnounceNexus)
	belowLayer.SetOnRevoke(p.OnRevoke)
}

// NewProviderTransactLayer creates a TransactLayer.
func NewProviderTransactLayer() *TransactLayer {
	return &TransactLayer{
		BaseLayer:         layer.NewBaseLayer(),
		NexusBySharedSeed: make(compat.NexusBySharedSeed),
	}
}

// AnnounceNexus announces a nexus.
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

// create a transact nexus.
func (p *TransactLayer) setupTransactNexus(belowNexus nexus.Nexus) *provider.TransactNexus {
	providerTransactNexus := provider.NewProviderTransactNexus(belowNexus, p)
	providerTransactNexus.HandleInvoiceRequest = p.handleInvoiceRequest
	providerTransactNexus.HandlePayRequest = p.handlePayRequest
	return providerTransactNexus
}

// RevokeNexus revokes a nexus.
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

// FulfilRequestInvoice pays a request invoice.
func (p *TransactLayer) FulfilRequestInvoice(nexusUUID, bolt11, requestReferenceUUID string) error {
	nexusID, _ := uuid.FromBytes([]byte(nexusUUID))
	if nx, ok := p.Nexuses.Get(nexusID); ok {
		providerTxNexus := nx.(*provider.TransactNexus)
		return providerTxNexus.NotifyInvoice(bolt11, requestReferenceUUID)
	}
	return nil
}

// NotifyPreImage notifies a preimage.
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

// NotifyProviderInfo notifies a provider info.
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

var _ layer.Base = &TransactLayer{}

// use an interface to call methods in the nexus.
var _ compat.ProviderTransactLayerInterface = &TransactLayer{}
