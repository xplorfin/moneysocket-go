package stack

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/provider"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/seller"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	nexusSeller "github.com/xplorfin/moneysocket-go/moneysocket/nexus/seller"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

type SellerStack struct {
	layer.BaseLayer
	websocketLayer              *websocket.OutgoingWebsocketLayer
	rendezvousLayer             *rendezvous.OutgoingRendezvousLayer
	providerLayer               *provider.Layer
	transactLayer               *provider.TransactLayer
	sellerLayer                 *seller.Layer
	handleProviderInfoRequest   compat.HandleProviderInfoRequest
	handleSellerInfoRequest     func() nexusSeller.Info
	handleInvoiceRequest        func(msats int64, requestUuid string)
	handleOpinionInvoiceRequest func(item string, requestUuid string)
	handlePayRequest            func(msats int64, requestUuid string)
	nexus                       nexus.Nexus
	sharedSeed                  beacon.SharedSeed
}

func NewSellerStack() *SellerStack {
	s := SellerStack{
		BaseLayer:       layer.NewBaseLayer(),
		websocketLayer:  websocket.NewOutgoingWebsocketLayer(),
		rendezvousLayer: rendezvous.NewOutgoingRendezvousLayer(),
		providerLayer:   provider.NewProviderLayer(),
		transactLayer:   provider.NewProviderTransactLayer(),
		sellerLayer:     seller.NewSellerLayer(),
	}
	s.SetupOutgoingWebsocketLayer()
	s.SetupOutgoingRendezvousLayer()
	s.SetupProviderLayer()
	s.SetupSellerLayer()

	return &s
}

func (s *SellerStack) SetupOutgoingWebsocketLayer() {
	s.websocketLayer.SetOnLayerEvent(func(layerName string, nexus nexus.Nexus, event string) {
		s.OnLayerEvent(message.OutgoingWebsocket, nexus, event)
	})
}

func (s *SellerStack) AnnounceNexus(belowNexus nexus.Nexus) {
	log.Println("provider stack got nexus")
	s.nexus = belowNexus
	s.sharedSeed = *belowNexus.SharedSeed()
	if s.OnAnnounce != nil {
		s.OnAnnounce(belowNexus)
	}
}

func (s *SellerStack) RevokeNexus(belowNexus nexus.Nexus) {
	log.Println("provider stack got nexus revoked")
	s.nexus = nil
	s.sharedSeed = beacon.NewSharedSeed()
	if s.OnAnnounce != nil {
		s.OnAnnounce(belowNexus)
	}
}

func (s *SellerStack) SellerNowReadyFromApp() {
	s.sellerLayer.SellerNowReadyFromApp()
}

func (s *SellerStack) SetupOutgoingRendezvousLayer() {
	s.rendezvousLayer.SetOnLayerEvent(func(layerName string, nexus nexus.Nexus, event string) {
		s.OnLayerEvent(message.OutgoingRendezvous, nexus, event)
	})
	s.rendezvousLayer.RegisterAboveLayer(s.websocketLayer)
}

func (s *SellerStack) SetOnStackEvent(onEvent func(layerName string, nexus nexus.Nexus, status string)) {
	s.SetOnLayerEvent(onEvent)
}

func (s *SellerStack) SetHandleProviderInfoRequest(hpr compat.HandleProviderInfoRequest) {
	s.handleProviderInfoRequest = hpr
}

func (s *SellerStack) SetupProviderLayer() {
	s.providerLayer.SetOnLayerEvent(func(layerName string, nexus nexus.Nexus, event string) {
		s.OnLayerEvent(message.Provider, nexus, event)
	})
	s.providerLayer.SetHandlerProvideInfoRequest(func(seed beacon.SharedSeed) account.Db {
		return s.handleProviderInfoRequest(seed)
	})
	s.providerLayer.RegisterAboveLayer(s.rendezvousLayer)
}

func (s *SellerStack) SetupProviderTransactLayer() {
	s.transactLayer.SetOnLayerEvent(func(layerName string, nexus nexus.Nexus, event string) {
		s.OnLayerEvent(message.ProviderTransact, nexus, event)
	})
	s.transactLayer.SetHandleInvoiceRequest(func(nexus nexus.Nexus, msats int64, requestUuid string) {
		s.HandleInvoiceRequest(msats, requestUuid)
	})
	s.transactLayer.RegisterAboveLayer(s.providerLayer)
}

func (s *SellerStack) ProviderNowReadyFromApp() {
	s.providerLayer.ProviderNowReadyFromApp()
}

func (s *SellerStack) SetHandleOpinionInvoiceRequest(handler func(item string, requestUuid string)) {
	s.handleOpinionInvoiceRequest = handler
}

func (s *SellerStack) SetupSellerLayer() {
	s.sellerLayer.SetOnLayerEvent(func(layerName string, nexus nexus.Nexus, event string) {
		s.OnLayerEvent(message.Seller, nexus, event)
	})
	s.sellerLayer.SetHandleOpinionInvoiceRequest(func(nx nexus.Nexus, itemId string, requestUuid string) {
		s.handleOpinionInvoiceRequest(itemId, requestUuid)
	})
	s.sellerLayer.SetHandleSellerInfoRequest(func() nexusSeller.Info {
		return s.handleSellerInfoRequest()
	})
	s.sellerLayer.SetOnAnnounce(func(nexus nexus.Nexus) {
		s.AnnounceNexus(nexus)
	})
	s.sellerLayer.SetOnRevoke(func(nexus nexus.Nexus) {
		s.OnRevoke(nexus)
	})

}

func (s *SellerStack) UpdatePrices() {
	s.nexus.(*nexusSeller.Nexus).UpdatePrices()
}

func (s *SellerStack) SetHandleSellerInfoRequest(handler func() nexusSeller.Info) {
	s.handleSellerInfoRequest = handler
}

func (s *SellerStack) HandleInvoiceRequest(msats int64, requestUUID string) {
	if s.handleInvoiceRequest != nil {
		s.handleInvoiceRequest(msats, requestUUID)
	}
}

func (s *SellerStack) HandleOpinionInvoiceRequest(itemID, requestUUID string) {
	if s.handleOpinionInvoiceRequest != nil {
		s.handleOpinionInvoiceRequest(itemID, requestUUID)
	}
}

func (s *SellerStack) SetHandleInvoiceRequest(handleInvoice func(msats int64, requestUuid string)) {
	s.handleInvoiceRequest = handleInvoice
}

func (s *SellerStack) SetHandlePayRequest(handlePayRequest func(msats int64, requestUuid string)) {
	s.handlePayRequest = handlePayRequest
}

func (s *SellerStack) DoDisconnect() {
	s.websocketLayer.InitiateCloseAll()
}
