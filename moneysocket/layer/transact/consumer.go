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

type ConsumerTransactLayer struct {
	layer.BaseLayer

	onInvoice             transact.OnInvoice
	onPreimage            transact.OnPreimage
	onProviderInfo        transact.OnProviderInfo
	consumerTransactNexus *transact.ConsumerTrackNexus
}

func NewConsumerTransactionLayer() *ConsumerTransactLayer {
	c := ConsumerTransactLayer{
		layer.NewBaseLayer(),
		nil,
		nil,
		nil,
		nil,
	}
	return &c
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (c *ConsumerTransactLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(c.AnnounceNexus)
	belowLayer.SetOnRevoke(c.OnRevoke)
}

// AnnounceNexus creates a new transact.ConsumerTransactNexus and registers it
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

// setup the consumer transaction enxus
func (c *ConsumerTransactLayer) SetupConsumerTransactionNexus(belowNexus nexus.Nexus) {
	ctn := transact.NewConsumerTransactNexus(belowNexus)
	c.consumerTransactNexus = ctn
	c.consumerTransactNexus.SetOnPreimage(c.onPreimage)
	c.consumerTransactNexus.SetOnInvoice(c.onInvoice)
	c.consumerTransactNexus.SetOnProviderInfo(c.onProviderInfo)
}

// call on invoice function
func (c *ConsumerTransactLayer) OnInvoice(transactNexus nexus.Nexus, invoice string, requestReferenceUuid string) {
	if c.onInvoice != nil {
		c.onInvoice(transactNexus, invoice, requestReferenceUuid)
	}
}

// set function to be called when on invoice is called
func (c *ConsumerTransactLayer) SetOnInvoice(invoice transact.OnInvoice) {
	c.onInvoice = invoice
}

// call on preimage function
func (c *ConsumerTransactLayer) OnPreImage(transactNexus nexus.Nexus, preimage string, requestReferenceUuid string) {
	if c.onPreimage != nil {
		c.onPreimage(transactNexus, preimage, requestReferenceUuid)
	}
}

// set function to be called when onPreImage is called
func (c *ConsumerTransactLayer) SetOnPreimage(preimage transact.OnPreimage) {
	c.onPreimage = preimage
}

func (c *ConsumerTransactLayer) OnProviderInfo(consumerTransactNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
	if c.onProviderInfo != nil {
		c.onProviderInfo(consumerTransactNexus, msg)
	}
}

func (c *ConsumerTransactLayer) SetOnProviderInfo(info transact.OnProviderInfo) {
	c.onProviderInfo = info
}

func (c *ConsumerTransactLayer) RequestPay(nexusUuid uuid.UUID, bolt11 string) (requestUuid uuid.UUID, err error) {
	if val, ok := c.Nexuses.Get(nexusUuid); !ok {
		consumerNexus := val.(*transact.ConsumerTrackNexus)
		res := consumerNexus.RequestPay(bolt11)
		return uuid.FromString(res)
	} else {
		return requestUuid, fmt.Errorf("nexus %s not online", nexusUuid)
	}
}

func (c *ConsumerTransactLayer) RequestInvoice(nexusUuid uuid.UUID, msats int64, description string) (requestUuid uuid.UUID, err error) {
	if val, ok := c.Nexuses.Get(nexusUuid); !ok {
		consumerNexus := val.(*transact.ConsumerTrackNexus)
		res := consumerNexus.RequestInvoice(msats, description)
		return uuid.FromString(res)
	} else {
		return requestUuid, fmt.Errorf("nexus %s not online", nexusUuid)
	}
}
