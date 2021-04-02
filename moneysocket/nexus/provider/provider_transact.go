package provider

import (
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	msg "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// TransactNexusName the name of the provider transaction
const TransactNexusName = "ProviderTransactNexus"

// TransactNexus is the provider transact nexus nam
type TransactNexus struct {
	*base.NexusBase
	HandleInvoiceRequest compat.HandleInvoiceRequest
	HandlePayRequest     compat.HandlePayRequest
}

// NewProviderTransactNexus creates a new TransactNexus
func NewProviderTransactNexus(belowNexus nexus.Nexus, layer layer.Base) *TransactNexus {
	nx := base.NewBaseNexusFull(TransactNexusName, belowNexus, layer)
	pn := TransactNexus{&nx, nil, nil}

	belowNexus.SetOnBinMessage(pn.OnBinMessage)
	belowNexus.SetOnMessage(pn.OnMessage)

	return &pn
}

// HandleLayerRequest handles a layer request
func (p *TransactNexus) HandleLayerRequest(req request.MoneysocketRequest) {
	if req.MessageType() == msg.InvoiceRequest {
		invoice := req.(request.Invoice)
		p.HandleInvoiceRequest(p, invoice.Msats, req.UUID())
	} else if req.MessageType() == msg.PayRequest {
		payRequest := req.(request.Pay)
		p.HandlePayRequest(p, payRequest.Bolt11, req.UUID())
	}
}

// IsLayerMessage checkks if a message can be processed by this layer
func (p *TransactNexus) IsLayerMessage(message msg.MoneysocketMessage) bool {
	if message.MessageClass() == msg.Request {
		return false
	}
	req := message.(request.MoneysocketRequest)
	return req.MessageType() == msg.InvoiceRequest || req.MessageType() == msg.PayRequest
}

// OnMessage processes a message
func (p *TransactNexus) OnMessage(belowNexus nexus.Nexus, message msg.MoneysocketMessage) {
	if !p.IsLayerMessage(message) {
		p.NexusBase.OnMessage(belowNexus, message)
	} else {
		req := message.(request.MoneysocketRequest)
		p.HandleLayerRequest(req)
	}
}

// OnBinMessage processes a  bin message
func (p *TransactNexus) OnBinMessage(baseNexus nexus.Nexus, msg []byte) {
	// pass
}

// NotifyInvoice notifies a new invoice
func (p *TransactNexus) NotifyInvoice(bolt11, requestReferenceUUID string) error {
	return p.Send(notification.NewNotifyInvoice(bolt11, requestReferenceUUID))
}

// NotifyPreimage notifies a preimage
func (p *TransactNexus) NotifyPreimage(preimage, requestReferenceUUID string) error {
	return p.Send(notification.NewNotifyPreimage(preimage, "", requestReferenceUUID))
}

// NotifyProviderInfo notifies a provider info
func (p *TransactNexus) NotifyProviderInfo(seed beacon.SharedSeed) error {
	pi := p.NexusBase.Layer.(compat.ProviderTransactLayerInterface)
	adb := pi.HandleProviderInfoRequest(seed)
	log.Println("notify provider")
	m := notification.NewNotifyProvider(adb.Details.AccountUUID, adb.Details.Payer(), adb.Details.Payee(), adb.Details.Wad, uuid.NewV4().String())
	return p.Send(m)
}

var _ nexus.Nexus = &TransactNexus{}
