package provider

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	message_base "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// NexusName is the ProviderNexus
const NexusName = "ProviderNexus"

// Nexus is the provider
type Nexus struct {
	*base.NexusBase
	// RequestReferenceUUID is the reference uuid
	RequestReferenceUUID      string
	handleProviderInfoRequest compat.HandleProviderInfoRequest
	// ProviderFinishedCb is the callback for when providers finish
	ProviderFinishedCb func(nx nexus.Nexus)
}

// NewProviderNexus creates a new nexus
func NewProviderNexus(belowNexus nexus.Nexus) *Nexus {
	baseNexus := base.NewBaseNexusBelow(NexusName, belowNexus)
	pn := Nexus{baseNexus, "", nil, nil}
	belowNexus.SetOnBinMessage(pn.OnBinMessage)
	belowNexus.SetOnMessage(pn.OnMessage)
	return &pn
}

// IsLayerMessage decides if a message should be processed by a layer
func (o *Nexus) IsLayerMessage(message message_base.MoneysocketMessage) bool {
	if message.MessageClass() == message_base.Request {
		return false
	}
	ntfn := message.(notification.MoneysocketNotification)
	return ntfn.RequestType() == message_base.ProviderRequest || ntfn.RequestType() == message_base.PingRequest
}

// NotifyProvider sends a notification to a provider
func (o *Nexus) NotifyProvider() {
	ss := o.SharedSeed()
	providerInfo := o.handleProviderInfoRequest(*ss)
	_ = o.Send(notification.NewNotifyProvider(o.UUID().String(), providerInfo.Details.Payer(), providerInfo.Details.Payee(), providerInfo.Details.Wad, o.RequestReferenceUUID))
}

// NotifyProviderNotReady notifies a provider is not ready yet
func (o *Nexus) NotifyProviderNotReady() {
	_ = o.Send(notification.NewNotifyProviderNotReady(o.RequestReferenceUUID))
}

// OnMessage processes the message for this layer
func (o *Nexus) OnMessage(belowNexus nexus.Nexus, msg message_base.MoneysocketMessage) {
	log.Println("provider nexus got message")
	if !o.IsLayerMessage(msg) {
		o.NexusBase.OnMessage(belowNexus, msg)
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

// NotifyPong sends a pong message
func (o *Nexus) NotifyPong() {
	_ = o.Send(notification.NewNotifyPong(o.RequestReferenceUUID))
}

// OnBinMessage calls the binary message handler
func (o *Nexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	log.Println("provider nexus got raw msg")
	o.NexusBase.OnBinMessage(belowNexus, msg)
}

// WaitForConsumer notifies the providerFinishedCb
func (o *Nexus) WaitForConsumer(providerFinishedCb func(nexus2 nexus.Nexus)) {
	o.ProviderFinishedCb = providerFinishedCb
}

// NotifyProviderReady notifies provider is ready
func (o *Nexus) NotifyProviderReady() {
	ss := o.SharedSeed()
	providerInfo := o.Layer.(compat.ProviderTransactLayerInterface).HandleProviderInfoRequest(*ss)
	if !providerInfo.Details.Ready() {
		panic("expected provider to be ready")
	}
	_ = o.Send(notification.NewNotifyProvider(o.UUID().String(), providerInfo.Details.Payer(), providerInfo.Details.Payee(), providerInfo.Details.Wad, o.RequestReferenceUUID))
}

// ProviderNowReady notifies provider is ready
func (o *Nexus) ProviderNowReady() {
	o.NotifyProviderReady()
}
