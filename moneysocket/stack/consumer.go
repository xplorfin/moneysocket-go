package stack

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/consumer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/transact"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	consumerNexus "github.com/xplorfin/moneysocket-go/moneysocket/nexus/consumer"
	nexus_transact "github.com/xplorfin/moneysocket-go/moneysocket/nexus/transact"
)

// ConsumerStack handles various events passed from/to child classes.
type ConsumerStack struct {
	layer.BaseLayer
	// gets called on layer registration
	sendStackEvent layer.OnLayerEventFn
	consumerLayer  *consumer.Layer
	transactLayer  *transact.ConsumerTransactLayer
	onPing         consumerNexus.OnPingFn
	// invoice event handler
	onInvoice nexus_transact.OnInvoice
	// preimage event handler
	onPreimage nexus_transact.OnPreimage
	// called on provder info
	onProviderInfo nexus_transact.OnProviderInfo

	nexus      nexusHelper.Nexus
	sharedSeed *beacon.SharedSeed
}

// NewConsumerStack creates a new ConsumerStack.
func NewConsumerStack() *ConsumerStack {
	c := ConsumerStack{}
	c.BaseLayer = layer.NewBaseLayer()
	return &c
}

// SetupConsumerLayer sets up a consumer layer.
func (c *ConsumerStack) SetupConsumerLayer(belowLayer layer.Base) {
	c.consumerLayer = consumer.NewConsumerLayer()
	c.consumerLayer.RegisterAboveLayer(belowLayer)
	c.consumerLayer.RegisterLayerEvent(c.sendStackEvent, message.Consumer)
	c.consumerLayer.SetOnPing(c.onPing)
}

// SetupTransactLayer sets up a transact layer.
func (c *ConsumerStack) SetupTransactLayer(belowLayer layer.Base) {
	c.transactLayer = transact.NewConsumerTransactLayer()
	c.transactLayer.RegisterAboveLayer(belowLayer)
	c.transactLayer.RegisterLayerEvent(c.sendStackEvent, message.Consumer)
	c.transactLayer.SetOnInvoice(c.onInvoice)
	c.transactLayer.SetOnPreimage(c.onPreimage)
	c.transactLayer.SetOnProviderInfo(c.onProviderInfo)
	c.transactLayer.SetOnAnnounce(func(nexus nexusHelper.Nexus) {
		c.AnnounceNexus(nexus)
	})
	c.transactLayer.SetOnRevoke(func(nexus nexusHelper.Nexus) {
		c.RevokeNexus(nexus)
	})
}

// OnPing sends a ping.
func (c *ConsumerStack) OnPing(transactNexus nexusHelper.Nexus, msecs int) {
	if c.onPing != nil {
		c.onPing(transactNexus, msecs)
	}
}

// SendStackEvent sends a stack event.
func (c *ConsumerStack) SendStackEvent(layerName string, nexus nexusHelper.Nexus, event string) {
	if c.sendStackEvent != nil {
		c.sendStackEvent(layerName, nexus, event)
	}
}

// OnInvoice calls the on invoice function.
func (c *ConsumerStack) OnInvoice(transactNexus nexusHelper.Nexus, invoice string, requestReferenceUUID string) {
	if c.onInvoice != nil {
		c.onInvoice(transactNexus, invoice, requestReferenceUUID)
	}
}

// SetOnInvoice sets function to be called when on invoice is called.
func (c *ConsumerStack) SetOnInvoice(invoice nexus_transact.OnInvoice) {
	c.onInvoice = invoice
}

// OnPreImage call on preimage function.
func (c *ConsumerStack) OnPreImage(transactNexus nexusHelper.Nexus, preimage, requestReferenceUUID string) {
	if c.onPreimage != nil {
		c.onPreimage(transactNexus, preimage, requestReferenceUUID)
	}
}

// SetOnPreimage sets function to be called when onPreImage is called.
func (c *ConsumerStack) SetOnPreimage(preimage nexus_transact.OnPreimage) {
	c.onPreimage = preimage
}

// SetOnPing sets function to be called when onPreImage is called.
func (c *ConsumerStack) SetOnPing(ping consumerNexus.OnPingFn) {
	c.onPing = ping
}

// SetOnProviderInfo sets function to be called when OnProviderInfo is called.
func (c *ConsumerStack) SetOnProviderInfo(info nexus_transact.OnProviderInfo) {
	c.onProviderInfo = info
}

// SetSendStackEvent sets a send stack callback.
func (c *ConsumerStack) SetSendStackEvent(handler layer.OnLayerEventFn) {
	c.sendStackEvent = handler
}

// OnProviderInfo handles a provider info event.
func (c *ConsumerStack) OnProviderInfo(consumerTransactNexus nexusHelper.Nexus, msg moneysocket_message.MoneysocketMessage) {
	if c.onProviderInfo != nil {
		c.onProviderInfo(consumerTransactNexus, msg)
	}
}

// AnnounceNexus announces a nexus.
func (c *ConsumerStack) AnnounceNexus(belowNexus nexusHelper.Nexus) {
	c.nexus = belowNexus
	c.sharedSeed = belowNexus.SharedSeed()
	if c.OnAnnounce != nil {
		c.OnAnnounce(belowNexus)
	}
}

// RevokeNexus removes the nexus from directories/layers. Calls OnRevoke.
func (c *ConsumerStack) RevokeNexus(belowNexus nexusHelper.Nexus) {
	c.nexus = nil
	c.sharedSeed = &beacon.SharedSeed{}
	if c.OnRevoke != nil {
		c.OnRevoke(belowNexus)
	}
}

// RequestInvoice requests an invoice for a given sat count.
func (c *ConsumerStack) RequestInvoice(msats int64, overrideRequestUUID, description string) {
	c.nexus.(compat.ConsumeNexusInterface).RequestInvoice(msats, overrideRequestUUID, description)
}

// RequestPay requests the lnd client to pay an invoice.
func (c *ConsumerStack) RequestPay(bolt11, overrideRequestUUID string) {
	c.nexus.(compat.ConsumeNexusInterface).RequestPay(bolt11, overrideRequestUUID)
}
