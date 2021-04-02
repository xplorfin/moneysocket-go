package transact

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/transact"
)

// ConsumerTransactLayer handles transaction.
type ConsumerTransactLayer struct {
	layer.BaseLayer

	onInvoice             transact.OnInvoice
	onPreimage            transact.OnPreimage
	onProviderInfo        transact.OnProviderInfo
	consumerTransactNexus *transact.ConsumerTrackNexus
}

// NewConsumerTransactLayer creates a ConsumerTransactLayer.
func NewConsumerTransactLayer() *ConsumerTransactLayer {
	c := ConsumerTransactLayer{
		layer.NewBaseLayer(),
		nil,
		nil,
		nil,
		nil,
	}
	return &c
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer.
func (c *ConsumerTransactLayer) RegisterAboveLayer(belowLayer layer.Base) {
	belowLayer.SetOnAnnounce(c.AnnounceNexus)
	belowLayer.SetOnRevoke(c.OnRevoke)
}

// AnnounceNexus creates a new transact.ConsumerTransactNexus and registers it.
func (c *ConsumerTransactLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	cn := transact.NewConsumerTransactNexus(belowNexus)
	c.consumerTransactNexus = cn
	c.TrackNexus(c.consumerTransactNexus, belowNexus)
	c.TrackNexusAnnounced(c.consumerTransactNexus)
	c.SendLayerEvent(c.consumerTransactNexus, message.NexusAnnounced)
	if c.OnAnnounce != nil {
		c.OnAnnounce(c.consumerTransactNexus)
	}
}

// SetupConsumerTransactionNexus registers the consumer.ConsumerTransactNexus.
func (c *ConsumerTransactLayer) SetupConsumerTransactionNexus(belowNexus nexus.Nexus) {
	ctn := transact.NewConsumerTransactNexus(belowNexus)
	c.consumerTransactNexus = ctn
	c.consumerTransactNexus.SetOnPreimage(c.onPreimage)
	c.consumerTransactNexus.SetOnInvoice(c.onInvoice)
	c.consumerTransactNexus.SetOnProviderInfo(c.onProviderInfo)
}

// OnInvoice calls an OnInvoice function.
func (c *ConsumerTransactLayer) OnInvoice(transactNexus nexus.Nexus, invoice string, requestReferenceUUID string) {
	if c.onInvoice != nil {
		c.onInvoice(transactNexus, invoice, requestReferenceUUID)
	}
}

// SetOnInvoice sets a function to be called when on invoice is called.
func (c *ConsumerTransactLayer) SetOnInvoice(invoice transact.OnInvoice) {
	c.onInvoice = invoice
}

// OnPreImage calls the registered on preimage function.
func (c *ConsumerTransactLayer) OnPreImage(transactNexus nexus.Nexus, preimage string, requestReferenceUUID string) {
	if c.onPreimage != nil {
		c.onPreimage(transactNexus, preimage, requestReferenceUUID)
	}
}

// SetOnPreimage sets a function to be called when onPreImage is called.
func (c *ConsumerTransactLayer) SetOnPreimage(preimage transact.OnPreimage) {
	c.onPreimage = preimage
}

// OnProviderInfo calls the registered ConsumerTransactLayer.OnProviderInfo.
func (c *ConsumerTransactLayer) OnProviderInfo(consumerTransactNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
	if c.onProviderInfo != nil {
		c.onProviderInfo(consumerTransactNexus, msg)
	}
}

// SetOnProviderInfo sets a provider info.
func (c *ConsumerTransactLayer) SetOnProviderInfo(info transact.OnProviderInfo) {
	c.onProviderInfo = info
}

// RequestPay requests payment for an invoice.
func (c *ConsumerTransactLayer) RequestPay(nexusUUID uuid.UUID, bolt11 string) (requestUUID uuid.UUID, err error) {
	if val, ok := c.Nexuses.Get(nexusUUID); !ok {
		consumerNexus := val.(*transact.ConsumerTrackNexus)
		res := consumerNexus.RequestPay(bolt11)
		return uuid.FromString(res)
	}

	return requestUUID, fmt.Errorf("nexus %s not online", nexusUUID)
}

// RequestInvoice requests an invoice.
func (c *ConsumerTransactLayer) RequestInvoice(nexusUUID uuid.UUID, msats int64, description string) (requestUUID uuid.UUID, err error) {
	if val, ok := c.Nexuses.Get(nexusUUID); !ok {
		consumerNexus := val.(*transact.ConsumerTrackNexus)
		res := consumerNexus.RequestInvoice(msats, description)
		return uuid.FromString(res)
	}
	return requestUUID, fmt.Errorf("nexus %s not online", nexusUUID)
}
