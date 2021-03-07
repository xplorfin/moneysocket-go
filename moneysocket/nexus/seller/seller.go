package seller

import (
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	msg "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

const SellerNexusName = "SellerNexus"

type SellerNexus struct {
	*base.BaseNexus
	handleSellerInfoRequest     func() SellerInfo
	sellerFinishedCb            func(nexus.Nexus)
	handleOpinionInvoiceRequest compat.HandleOpinionInvoiceRequest
}

// seller info message
type SellerInfo struct {
	// wether or not the seller is ready
	Ready bool `json:"ready"`
	// wether or not the uuid works
	SellerUUID string              `json:"seller_uuid"`
	Items      []notification.Item `json:"items"`
}

func NewSellerNexus(belowNexus nexus.Nexus, layer layer.Layer) SellerNexus {
	baseNexus := base.NewBaseNexusFull(SellerNexusName, belowNexus, layer)
	sn := SellerNexus{
		&baseNexus,
		nil,
		nil,
		nil,
	}
	belowNexus.SetOnBinMessage(sn.OnBinMessage)
	belowNexus.SetOnMessage(sn.OnMessage)
	return sn
}

func (s *SellerNexus) IsLayerMessage(message msg.MoneysocketMessage) bool {
	if message.MessageClass() == msg.Request {
		return false
	}

	req := message.(request.MoneysocketRequest)
	return req.MessageType() == msg.RequestOpinionSeller || req.MessageType() == msg.RequestOpinionInvoice
}

func (s *SellerNexus) notifySeller(requestReferenceUuid string) error {
	sellerInfo := s.handleSellerInfoRequest()
	return s.Send(notification.NewNotifyOpinionSeller(sellerInfo.SellerUUID, sellerInfo.Items, requestReferenceUuid))
}

func (s *SellerNexus) UpdatePrices() {
	s.notifySeller(uuid.NewV4().String())
}

func (s *SellerNexus) SetHandleOpinionInvoiceRequest(invoiceRequest compat.HandleOpinionInvoiceRequest) {
	s.handleOpinionInvoiceRequest = invoiceRequest
}

func (s *SellerNexus) SetHandleSellerInfoRequest(handler func() SellerInfo) {
	s.handleSellerInfoRequest = handler
}

func (s *SellerNexus) notifySellerNotReady(requestReferenceUuid string) error {
	return s.Send(notification.NewNotifyOpinionSellerNotReady(requestReferenceUuid))
}

func (s *SellerNexus) OnMessage(baseNexus nexus.Nexus, message msg.MoneysocketMessage) {
	log.Println("provider nexus got message from below")
	if !s.IsLayerMessage(message) {
		s.BaseNexus.OnMessage(baseNexus, message)
	}
	// message request
	nx := message.(notification.MoneysocketNotification)
	if nx.RequestType() == msg.RequestOpinionSeller {
		sharedSeed := s.SharedSeed()
		sellerInfo := s.handleSellerInfoRequest()
		if sellerInfo.Ready {
			_ = s.notifySeller(nx.RequestReferenceUuid())
			s.sellerFinishedCb(s)
		} else {
			l := s.Layer.(compat.SellingLayerInterface)
			_ = s.notifySellerNotReady(nx.RequestReferenceUuid())
			l.NexusWaitingForApp(sharedSeed, s)
		}
	} else if nx.RequestType() == msg.RequestOpinionInvoice {
		mg := message.(request.RequestOpinionInvoice)
		s.handleOpinionInvoiceRequest(s, mg.ItemId, nx.RequestReferenceUuid())
	}
}

func (s *SellerNexus) NotifySellerNotReady(requestReferenceUuid string) error {
	return s.Send(notification.NewNotifyOpinionSellerNotReady(requestReferenceUuid))
}

func (s *SellerNexus) OnBinMessage(belowNexus nexus.Nexus, msgBytes []byte) {
	s.BaseNexus.OnBinMessage(belowNexus, msgBytes)
}

func (s *SellerNexus) NotifyOpinionInvoice(bolt11, requestReferenceUuid string) error {
	return s.Send(notification.NewNotifyOpinionInvoice(requestReferenceUuid, bolt11))
}

func (s *SellerNexus) WaitForBuyer(sellerFinishedCb func(nexus.Nexus)) {
	s.sellerFinishedCb = sellerFinishedCb
}

func (s *SellerNexus) SellerNowReady() {
	s.notifySeller(uuid.NewV4().String())
	s.sellerFinishedCb(s)
}
