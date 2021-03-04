package terminus

import (
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	message_base "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

const TerminusNexusName = "TerminusNexus"

type TerminusNexus struct {
	*base.BaseNexus
	handleInvoiceRequest      compat.HandleInvoiceRequest
	handlePayRequest          compat.HandlePayRequest
	handleProviderInfoRequest compat.HandleProviderInfoRequest
}

func (o TerminusNexus) IsLayerMessage(message message_base.MoneysocketMessage) bool {
	if message.MessageClass() != message_base.Request {
		return false
	}
	req := message.(request.MoneysocketRequest)
	return req.MessageType() == message_base.PayRequest || req.MessageType() == message_base.InvoiceRequest
}

func (o *TerminusNexus) OnMessage(belowNexus nexus.Nexus, message message_base.MoneysocketMessage) {
	log.Println("terminus nexus got msg")
	if !o.IsLayerMessage(message) {
		o.BaseNexus.OnMessage(belowNexus, message)
		return
	}
	req := message.(request.MoneysocketRequest)
	if req.MessageType() == message_base.PayRequest {
		payReq := req.(request.RequestPay)
		o.handlePayRequest(o, payReq.Bolt11, payReq.Uuid())
	}
	if req.MessageType() == message_base.InvoiceRequest {
		iReq := req.(request.RequestInvoice)
		o.handleInvoiceRequest(o, iReq.Msats, iReq.Uuid())
		// TODO we need a notify invoice here
		panic("method not yet implemented")
	}
}

func (o *TerminusNexus) OnBinMessage(belowNexus nexus.Nexus, msgBytes []byte) {
	log.Println("terminus nexus got raw msg")
	o.BaseNexus.OnBinMessage(belowNexus, msgBytes)
}

func (o *TerminusNexus) NotifyPreimage(preimage, requestReferenceUuid string) {
	_ = o.Send(notification.NewNotifyPreimage(preimage, "", requestReferenceUuid))
}

func (o *TerminusNexus) NotifyProviderInfo(ss beacon.SharedSeed) {
	pi := o.handleProviderInfoRequest(ss)
	m := notification.NewNotifyProvider(pi.Details.AccountUuid, pi.Details.Payer(), pi.Details.Payee(), pi.Details.Wad, uuid.NewV4().String())
	_ = o.Send(m)
}

func NewTerminusNexus(below nexus.Nexus, layer layer.Layer) TerminusNexus {
	bn := base.NewBaseNexusFull(TerminusNexusName, below, layer)
	return TerminusNexus{
		BaseNexus:                 &bn,
		handleInvoiceRequest:      nil,
		handlePayRequest:          nil,
		handleProviderInfoRequest: nil,
	}
}

var _ nexus.Nexus = &TerminusNexus{}
