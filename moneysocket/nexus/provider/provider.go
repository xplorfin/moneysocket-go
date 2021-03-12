package provider

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	message_base "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

const ProviderNexusName = "ProviderNexus"

type ProviderNexus struct {
	*base.BaseNexus
	RequestReferenceUuid      string
	handleProviderInfoRequest compat.HandleProviderInfoRequest
	ProviderFinishedCb        func(nx nexus.Nexus)
}

func NewProviderNexus(belowNexus nexus.Nexus) *ProviderNexus {
	baseNexus := base.NewBaseNexusBelow(ProviderNexusName, belowNexus)
	pn := ProviderNexus{baseNexus, "", nil, nil}
	belowNexus.SetOnBinMessage(pn.OnBinMessage)
	belowNexus.SetOnMessage(pn.OnMessage)
	return &pn
}

func (o *ProviderNexus) IsLayerMessage(message message_base.MoneysocketMessage) bool {
	if message.MessageClass() == message_base.Request {
		return false
	}
	ntfn := message.(notification.MoneysocketNotification)
	return ntfn.RequestType() == message_base.ProviderRequest || ntfn.RequestType() == message_base.PingRequest
}

func (o *ProviderNexus) NotifyProvider() {
	ss := o.SharedSeed()
	providerInfo := o.handleProviderInfoRequest(*ss)
	o.Send(notification.NewNotifyProvider(o.Uuid().String(), providerInfo.Details.Payer(), providerInfo.Details.Payee(), providerInfo.Details.Wad, o.RequestReferenceUuid))
}

func (o *ProviderNexus) NotifyProviderNotReady() {
	o.Send(notification.NewNotifyProviderNotReady(o.RequestReferenceUuid))
}

func (o *ProviderNexus) OnMessage(belowNexus nexus.Nexus, msg message_base.MoneysocketMessage) {
	log.Println("provider nexus got message")
	if !o.IsLayerMessage(msg) {
		o.BaseNexus.OnMessage(belowNexus, msg)
		return
	}
	ntfn := msg.(notification.MoneysocketNotification)
	if ntfn.RequestType() == message_base.ProviderRequest {
		ss := belowNexus.SharedSeed()
		providerInfo := o.handleProviderInfoRequest(*ss)
		if providerInfo.Details.Ready() {
			o.NotifyProvider()
			o.ProviderFinishedCb(o)
		} else {
			o.NotifyProviderNotReady()
			o.Layer.(compat.SellingLayerInterface).NexusWaitingForApp(ss, o)
		}
	} else if ntfn.RequestType() == message_base.PingRequest {
		o.NotifyPong()
	}

}

func (o *ProviderNexus) NotifyPong() {
	o.Send(notification.NewNotifyPong(o.RequestReferenceUuid))
}

func (o *ProviderNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	log.Println("provider nexus got raw msg")
	o.BaseNexus.OnBinMessage(belowNexus, msg)
}

func (o *ProviderNexus) WaitForConsumer(providerFinishedCb func(nexus2 nexus.Nexus)) {
	o.ProviderFinishedCb = providerFinishedCb
}

func (o *ProviderNexus) NotifyProviderReady() {
	ss := o.SharedSeed()
	providerInfo := o.Layer.(compat.ProviderTransactLayerInterface).HandleProviderInfoRequest(*ss)
	if !providerInfo.Details.Ready() {
		panic("expected provider to be ready")
	}
	o.Send(notification.NewNotifyProvider(o.Uuid().String(), providerInfo.Details.Payer(), providerInfo.Details.Payee(), providerInfo.Details.Wad, o.RequestReferenceUuid))
}
func (o *ProviderNexus) ProviderNowReady() {
	o.NotifyProviderReady()
}
