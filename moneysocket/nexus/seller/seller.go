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

// SellerNexusName is the seller name.
const SellerNexusName = "SellerNexus"

// Nexus - SellerNexus handler.
type Nexus struct {
	*base.NexusBase
	handleSellerInfoRequest     func() Info
	sellerFinishedCb            func(nexus.Nexus)
	handleOpinionInvoiceRequest compat.HandleOpinionInvoiceRequest
}

// Info seller message.
type Info struct {
	// wether or not the seller is ready
	Ready bool `json:"ready"`
	// wether or not the uuid works
	SellerUUID string              `json:"seller_uuid"`
	Items      []notification.Item `json:"items"`
}

// NewSellerNexus creates a nexus.
func NewSellerNexus(belowNexus nexus.Nexus, layer layer.Base) *Nexus {
	baseNexus := base.NewBaseNexusFull(SellerNexusName, belowNexus, layer)
	sn := Nexus{
		&baseNexus,
		nil,
		nil,
		nil,
	}
	belowNexus.SetOnBinMessage(sn.OnBinMessage)
	belowNexus.SetOnMessage(sn.OnMessage)
	return &sn
}

// IsLayerMessage processes a layer message.
func (s *Nexus) IsLayerMessage(message msg.MoneysocketMessage) bool {
	if message.MessageClass() == msg.Request {
		return false
	}

	req := message.(request.MoneysocketRequest)
	return req.MessageType() == msg.RequestOpinionSeller || req.MessageType() == msg.RequestOpinionInvoice
}

func (s *Nexus) notifySeller(requestReferenceUUID string) error {
	sellerInfo := s.handleSellerInfoRequest()
	return s.Send(notification.NewNotifyOpinionSeller(sellerInfo.SellerUUID, sellerInfo.Items, requestReferenceUUID))
}

// UpdatePrices updates seller prices.
func (s *Nexus) UpdatePrices() {
	_ = s.notifySeller(uuid.NewV4().String())
}

// SetHandleOpinionInvoiceRequest sets an invoice request.
func (s *Nexus) SetHandleOpinionInvoiceRequest(invoiceRequest compat.HandleOpinionInvoiceRequest) {
	s.handleOpinionInvoiceRequest = invoiceRequest
}

// SetHandleSellerInfoRequest sets the info request.
func (s *Nexus) SetHandleSellerInfoRequest(handler func() Info) {
	s.handleSellerInfoRequest = handler
}

func (s *Nexus) notifySellerNotReady(requestReferenceUUID string) error {
	return s.Send(notification.NewNotifyOpinionSellerNotReady(requestReferenceUUID))
}

// OnMessage processes the message.
func (s *Nexus) OnMessage(baseNexus nexus.Nexus, message msg.MoneysocketMessage) {
	log.Println("provider nexus got message from below")
	if !s.IsLayerMessage(message) {
		s.NexusBase.OnMessage(baseNexus, message)
	}
	// message request
	nx := message.(notification.MoneysocketNotification)
	if nx.RequestType() == msg.RequestOpinionSeller {
		sharedSeed := s.SharedSeed()
		sellerInfo := s.handleSellerInfoRequest()
		if sellerInfo.Ready {
			_ = s.notifySeller(nx.RequestReferenceUUID())
			s.sellerFinishedCb(s)
		} else {
			l := s.Layer.(compat.SellingLayerInterface)
			_ = s.notifySellerNotReady(nx.RequestReferenceUUID())
			l.NexusWaitingForApp(sharedSeed, s)
		}
	} else if nx.RequestType() == msg.RequestOpinionInvoice {
		mg := message.(request.OpinionInvoice)
		s.handleOpinionInvoiceRequest(s, mg.ItemID, nx.RequestReferenceUUID())
	}
}

// NotifySellerNotReady notifies the seller is not ready.
func (s *Nexus) NotifySellerNotReady(requestReferenceUUID string) error {
	return s.Send(notification.NewNotifyOpinionSellerNotReady(requestReferenceUUID))
}

// OnBinMessage calls the on bin message callback.
func (s *Nexus) OnBinMessage(belowNexus nexus.Nexus, msgBytes []byte) {
	s.NexusBase.OnBinMessage(belowNexus, msgBytes)
}

// NotifyOpinionInvoice notifies an invoice.
func (s *Nexus) NotifyOpinionInvoice(bolt11, requestReferenceUUID string) error {
	return s.Send(notification.NewNotifyOpinionInvoice(requestReferenceUUID, bolt11))
}

// WaitForBuyer initializes the sellerFinishedCb.
func (s *Nexus) WaitForBuyer(sellerFinishedCb func(nexus.Nexus)) {
	s.sellerFinishedCb = sellerFinishedCb
}

// SellerNowReady notifies the seller is now ready.
func (s *Nexus) SellerNowReady() {
	_ = s.notifySeller(uuid.NewV4().String())
	s.sellerFinishedCb(s)
}
