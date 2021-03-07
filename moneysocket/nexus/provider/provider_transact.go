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

const ProviderTransactNexusName = "ProviderTransactNexus"

type ProviderTransactNexus struct {
	*base.BaseNexus
	HandleInvoiceRequest compat.HandleInvoiceRequest
	HandlePayRequest     compat.HandlePayRequest
}

func NewProviderTransactNexus(belowNexus nexus.Nexus, layer layer.Layer) ProviderTransactNexus {
	nx := base.NewBaseNexusFull(ProviderTransactNexusName, belowNexus, layer)
	belowNexus.SetOnBinMessage(nx.OnBinMessage)
	belowNexus.SetOnMessage(nx.OnMessage)

	return ProviderTransactNexus{&nx, nil, nil}
}

// handle layer request
func (p *ProviderTransactNexus) HandleLayerRequest(req request.MoneysocketRequest) {
	if req.MessageType() == msg.InvoiceRequest {
		invoice := req.(request.RequestInvoice)
		p.HandleInvoiceRequest(p, invoice.Msats, req.Uuid())
	} else if req.MessageType() == msg.PayRequest {
		payRequest := req.(request.RequestPay)
		p.HandlePayRequest(p, payRequest.Bolt11, req.Uuid())
	}
}

func (p *ProviderTransactNexus) IsLayerMessage(message msg.MoneysocketMessage) bool {
	if message.MessageClass() == msg.Request {
		return false
	}
	req := message.(request.MoneysocketRequest)
	return req.MessageType() == msg.InvoiceRequest || req.MessageType() == msg.PayRequest
}

func (p *ProviderTransactNexus) OnMessage(belowNexus nexus.Nexus, message msg.MoneysocketMessage) {
	if !p.IsLayerMessage(message) {
		p.BaseNexus.OnMessage(belowNexus, message)
	} else {
		req := message.(request.MoneysocketRequest)
		p.HandleLayerRequest(req)
	}
}

func (p *ProviderTransactNexus) OnBinMessage(baseNexus nexus.Nexus, msg []byte) {
	// pass
}

func (p *ProviderTransactNexus) NotifyInvoice(bolt11, requestReferenceUuid string) error {
	return p.Send(notification.NewNotifyInvoice(bolt11, requestReferenceUuid))
}

func (p *ProviderTransactNexus) NotifyPreimage(preimage, requestReferenceUuid string) error {
	return p.Send(notification.NewNotifyPreimage(preimage, "", requestReferenceUuid))
}

func (p *ProviderTransactNexus) NotifyProviderInfo(seed beacon.SharedSeed) error {
	pi := p.BaseNexus.Layer.(compat.ProviderTransactLayerInterface)
	adb := pi.HandleProviderInfoRequest(seed)
	log.Println("notify provider")
	m := notification.NewNotifyProvider(adb.Details.AccountUuid, adb.Details.Payer(), adb.Details.Payee(), adb.Details.Wad, uuid.NewV4().String())
	return p.Send(m)
}

var _ nexus.Nexus = &ProviderTransactNexus{}
