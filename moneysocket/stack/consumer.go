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

// handles various events passed from/to child classes
type ConsumerStack struct {
	layer.BaseLayer
	// gets called on layer registration
	sendStackEvent layer.OnLayerEventFn
	consumerLayer  *consumer.ConsumerLayer
	transactLayer  *transact.ConsumerTransactLayer
	onPing         consumerNexus.OnPingFn
	// invoice event handler
	onInvoice nexus_transact.OnInvoice
	// preimage event handler
	onPreimage nexus_transact.OnPreimage
	// called on provder info
	onProviderInfo nexus_transact.OnProviderInfo

	nexus      nexusHelper.Nexus
	sharedSeed beacon.SharedSeed
}

func NewConsumerStack() *ConsumerStack {
	c := ConsumerStack{}
	c.BaseLayer = layer.NewBaseLayer()
	return &c
}

func (c *ConsumerStack) SetupConsumerLayer(belowLayer layer.Layer) {
	c.consumerLayer = consumer.NewConsumerLayer()
	c.consumerLayer.RegisterAboveLayer(belowLayer)
	c.consumerLayer.RegisterLayerEvent(c.sendStackEvent, message.Consumer)
	c.consumerLayer.SetOnPing(c.onPing)
}

func (c *ConsumerStack) SetupTransactLayer(belowLayer layer.Layer) {
	c.transactLayer = transact.NewConsumerTransactionLayer()
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

func (c *ConsumerStack) OnPing(transactNexus nexusHelper.Nexus, msecs int) {
	if c.onPing != nil {
		c.onPing(transactNexus, msecs)
	}
}

func (c *ConsumerStack) SendStackEvent(layerName string, nexus nexusHelper.Nexus, event string) {
	if c.sendStackEvent != nil {
		c.sendStackEvent(layerName, nexus, event)
	}
}

// call on invoice function
func (c *ConsumerStack) OnInvoice(transactNexus nexusHelper.Nexus, invoice string, requestReferenceUuid string) {
	if c.onInvoice != nil {
		c.onInvoice(transactNexus, invoice, requestReferenceUuid)
	}
}

// set function to be called when on invoice is called
func (c *ConsumerStack) SetOnInvoice(invoice nexus_transact.OnInvoice) {
	c.onInvoice = invoice
}

// call on preimage function
func (c *ConsumerStack) OnPreImage(transactNexus nexusHelper.Nexus, preimage string, requestReferenceUuid string) {
	if c.onPreimage != nil {
		c.onPreimage(transactNexus, preimage, requestReferenceUuid)
	}
}

// set function to be called when onPreImage is called
func (c *ConsumerStack) SetOnPreimage(preimage nexus_transact.OnPreimage) {
	c.onPreimage = preimage
}

// set function to be called when onPreImage is called
func (c *ConsumerStack) SetOnPing(ping consumerNexus.OnPingFn) {
	c.onPing = ping
}

// set function to be called when OnProviderInfo is called
func (c *ConsumerStack) SetOnProviderInfo(info nexus_transact.OnProviderInfo) {
	c.onProviderInfo = info
}

func (c *ConsumerStack) SetSendStackEvent(handler layer.OnLayerEventFn) {
	c.sendStackEvent = handler
}

func (c *ConsumerStack) OnProviderInfo(consumerTransactNexus nexusHelper.Nexus, msg moneysocket_message.MoneysocketMessage) {
	if c.onProviderInfo != nil {
		c.onProviderInfo(consumerTransactNexus, msg)
	}
}

func (c *ConsumerStack) AnnounceNexus(belowNexus nexusHelper.Nexus) {
	c.nexus = belowNexus
	c.sharedSeed = belowNexus.SharedSeed()
	if c.OnAnnounce != nil {
		c.OnAnnounce(belowNexus)
	}
}

func (c *ConsumerStack) RevokeNexus(belowNexus nexusHelper.Nexus) {
	c.nexus = nil
	c.sharedSeed = beacon.SharedSeed{}
	if c.OnRevoke != nil {
		c.OnRevoke(belowNexus)
	}
}

func (c *ConsumerStack) RequestInvoice(msats int64, overrideRequestUuid, description string) {
	c.nexus.(compat.ConsumeNexusInterface).RequestInvoice(msats, overrideRequestUuid, description)
}
func (c *ConsumerStack) RequestPay(bolt11, overrideRequestUuid string) {
	c.nexus.(compat.ConsumeNexusInterface).RequestPay(bolt11, overrideRequestUuid)
}
