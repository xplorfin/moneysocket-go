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

type SellerLayer struct {
	layer.BaseLayer
	// app's we're waiting for
	WaitingForApp compat.WaitingForApp
	// nexus by shared seed
	NexusBySharedSeed compat.NexusBySharedSeed

	handleOpinionInvoiceRequest compat.HandleOpinionInvoiceRequest

	handleSellerInfoRequest func() seller.SellerInfo

	SellerNexus seller.SellerNexus
}

func NewSellerLayer() *SellerLayer {
	return &SellerLayer{
		BaseLayer:         layer.NewBaseLayer(),
		WaitingForApp:     make(map[string]nexus.Nexus),
		NexusBySharedSeed: make(compat.NexusBySharedSeed),
	}
}

func (s *SellerLayer) HandleOpinionInvoiceRequest(nx nexus.Nexus, itemId string, requestUuid string) {
	s.handleOpinionInvoiceRequest(nx, itemId, requestUuid)
}

func (s *SellerLayer) SetHandleOpinionInvoiceRequest(request compat.HandleOpinionInvoiceRequest) {
	s.handleOpinionInvoiceRequest = request
}

func (s *SellerLayer) SetHandleSellerInfoRequest(handler func() seller.SellerInfo) {
	s.handleSellerInfoRequest = handler
}

func (s *SellerLayer) SetupSellerNexus(belowNexus nexus.Nexus) seller.SellerNexus {
	n := seller.NewSellerNexus(belowNexus, s)
	n.SetHandleOpinionInvoiceRequest(func(nx nexus.Nexus, itemId string, requestUuid string) {
		s.handleOpinionInvoiceRequest(nx, itemId, requestUuid)
	})
	n.SetHandleSellerInfoRequest(func() seller.SellerInfo {
		return s.handleSellerInfoRequest()
	})
	return n
}

func (s *SellerLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	s.SetOnAnnounce(belowLayer.AnnounceNexus)
	s.SetOnRevoke(belowLayer.RevokeNexus)
}

func (s *SellerLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	log.Println("buyer layer got nexus, starting handshake")
	s.SellerNexus = s.SetupSellerNexus(belowNexus)
	s.TrackNexus(&s.SellerNexus, belowNexus)

	s.SendLayerEvent(&s.SellerNexus, message.NexusWaiting)
	s.SellerNexus.WaitForBuyer(s.sellerFinishedCb)
}

func (s *SellerLayer) sellerFinishedCb(nx nexus.Nexus) {
	s.TrackNexusAnnounced(nx)
	s.SendLayerEvent(nx, message.NexusAnnounced)
	if s.OnAnnounce != nil {
		s.OnAnnounce(nx)
	}
}

func (s *SellerLayer) SellerNowReadyFromApp() {
	log.Println("-- seller now ready")
	for seed, _ := range s.WaitingForApp { // nolint
		log.Println("-- unwaiting for app")
		sellerNexus := s.WaitingForApp[seed]
		delete(s.WaitingForApp, seed)
		sellerNexus.(*seller.SellerNexus).SellerNowReady()
	}
}

func (s *SellerLayer) NexusWaitingForApp(seed *beacon.SharedSeed, sellerNexus nexus.Nexus) {
	log.Println("-- waiting for app")
	s.WaitingForApp[seed.ToString()] = sellerNexus
}

var _ compat.SellingLayerInterface = &SellerLayer{}
var _ layer.Layer = &SellerLayer{}
