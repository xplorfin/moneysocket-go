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

const NexusName = "TerminusNexus"

type Nexus struct {
	*base.NexusBase
	handleInvoiceRequest      compat.HandleInvoiceRequest
	handlePayRequest          compat.HandlePayRequest
	handleProviderInfoRequest compat.HandleProviderInfoRequest
}

func (o Nexus) IsLayerMessage(message message_base.MoneysocketMessage) bool {
	if message.MessageClass() != message_base.Request {
		return false
	}
	req := message.(request.MoneysocketRequest)
	return req.MessageType() == message_base.PayRequest || req.MessageType() == message_base.InvoiceRequest
}

func (o *Nexus) OnMessage(belowNexus nexus.Nexus, message message_base.MoneysocketMessage) {
	log.Println("terminus nexus got msg")
	if !o.IsLayerMessage(message) {
		o.NexusBase.OnMessage(belowNexus, message)
		return
	}
	req := message.(request.MoneysocketRequest)
	if req.MessageType() == message_base.PayRequest {
		payReq := req.(request.Pay)
		o.handlePayRequest(o, payReq.Bolt11, payReq.UUID())
	}
	if req.MessageType() == message_base.InvoiceRequest {
		iReq := req.(request.Invoice)
		o.handleInvoiceRequest(o, iReq.Msats, iReq.UUID())
		// TODO we need a notify invoice here
		panic("method not yet implemented")
	}
}

func (o *Nexus) OnBinMessage(belowNexus nexus.Nexus, msgBytes []byte) {
	log.Println("terminus nexus got raw msg")
	o.NexusBase.OnBinMessage(belowNexus, msgBytes)
}

func (o *Nexus) NotifyPreimage(preimage, requestReferenceUUID string) {
	_ = o.Send(notification.NewNotifyPreimage(preimage, "", requestReferenceUUID))
}

func (o *Nexus) NotifyProviderInfo(ss beacon.SharedSeed) {
	pi := o.handleProviderInfoRequest(ss)
	m := notification.NewNotifyProvider(pi.Details.AccountUUID, pi.Details.Payer(), pi.Details.Payee(), pi.Details.Wad, uuid.NewV4().String())
	_ = o.Send(m)
}

func NewTerminusNexus(below nexus.Nexus, layer layer.LayerBase) Nexus {
	bn := base.NewBaseNexusFull(NexusName, below, layer)
	return Nexus{
		NexusBase:                 &bn,
		handleInvoiceRequest:      nil,
		handlePayRequest:          nil,
		handleProviderInfoRequest: nil,
	}
}

var _ nexus.Nexus = &Nexus{}
