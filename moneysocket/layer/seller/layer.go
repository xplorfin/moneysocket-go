package seller

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/seller"
)

// SellerLayer is meant to simulate the seller layer in the js seller app
// this struct should not be initialized directly, the NewSellerLayer() method
// below should be used instead
type Layer struct {
	layer.BaseLayer
	// nexuses's we're waiting to initialize
	WaitingForApp compat.WaitingForApp
	// nexus list by shared seed
	NexusBySharedSeed compat.NexusBySharedSeed
	// event handler for an invoice request (supplied by client)
	handleOpinionInvoiceRequest compat.HandleOpinionInvoiceRequest
	// event handler for an info request (supplied by client)
	handleSellerInfoRequest func() seller.Info
	// the seller nexus object
	SellerNexus *seller.Nexus
}

// Create a new seller layer
func NewSellerLayer() *Layer {
	return &Layer{
		BaseLayer:         layer.NewBaseLayer(),
		WaitingForApp:     make(map[string]nexus.Nexus),
		NexusBySharedSeed: make(compat.NexusBySharedSeed),
	}
}

// Calls the client supplied handleOpinionInvoiceRequest method if present
func (s *Layer) HandleOpinionInvoiceRequest(nx nexus.Nexus, itemID string, requestUUID string) {
	s.handleOpinionInvoiceRequest(nx, itemID, requestUUID)
}

// Sets the client supplied handleOpinionInvoiceRequest method. This method is null by default
func (s *Layer) SetHandleOpinionInvoiceRequest(request compat.HandleOpinionInvoiceRequest) {
	s.handleOpinionInvoiceRequest = request
}

// Sets the client supplied handleSellerInfoRequest method. This method is null by default
func (s *Layer) SetHandleSellerInfoRequest(handler func() seller.Info) {
	s.handleSellerInfoRequest = handler
}

// Starts the seller nexus and initializes the callbacks
func (s *Layer) SetupSellerNexus(belowNexus nexus.Nexus) *seller.Nexus {
	n := seller.NewSellerNexus(belowNexus, s)
	n.SetHandleOpinionInvoiceRequest(func(nx nexus.Nexus, itemId string, requestUuid string) {
		s.handleOpinionInvoiceRequest(nx, itemId, requestUuid)
	})
	n.SetHandleSellerInfoRequest(func() seller.Info {
		return s.handleSellerInfoRequest()
	})
	return n
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (s *Layer) RegisterAboveLayer(belowLayer layer.LayerBase) {
	s.SetOnAnnounce(belowLayer.AnnounceNexus)
	s.SetOnRevoke(belowLayer.RevokeNexus)
}

// AnnounceNexus creates a new SellerNexus and registers it
// also registers the sellerFinished callback
func (s *Layer) AnnounceNexus(belowNexus nexus.Nexus) {
	log.Println("buyer layer got nexus, starting handshake")
	s.SellerNexus = s.SetupSellerNexus(belowNexus)
	s.TrackNexus(s.SellerNexus, belowNexus)

	s.SendLayerEvent(s.SellerNexus, message.NexusWaiting)
	s.SellerNexus.WaitForBuyer(s.sellerFinishedCb)
}

// callback for seller finished
func (s *Layer) sellerFinishedCb(nx nexus.Nexus) {
	s.TrackNexusAnnounced(nx)
	s.SendLayerEvent(nx, message.NexusAnnounced)
	if s.OnAnnounce != nil {
		s.OnAnnounce(nx)
	}
}

// sets the seller's status to ready (open store)
func (s *Layer) SellerNowReadyFromApp() {
	log.Println("-- seller now ready")
	for seed, _ := range s.WaitingForApp { // nolint
		log.Println("-- unwaiting for app")
		sellerNexus := s.WaitingForApp[seed]
		delete(s.WaitingForApp, seed)
		sellerNexus.(*seller.Nexus).SellerNowReady()
	}
}

// sets the seller's status to waiting
func (s *Layer) NexusWaitingForApp(seed *beacon.SharedSeed, sellerNexus nexus.Nexus) {
	log.Println("-- waiting for app")
	s.WaitingForApp[seed.ToString()] = sellerNexus
}

// make sure that the seller layer is compatible with interfaces used for calling
// non-standard layer methods and standard layer methods
var _ compat.SellingLayerInterface = &Layer{}
var _ layer.LayerBase = &Layer{}
