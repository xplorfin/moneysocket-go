package transact

import (
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// OnInvoice handles an invoice request
type OnInvoice func(transactNexus nexus.Nexus, invoice string, requestReferenceUuid string)

// OnPreimage handles a preimage request
type OnPreimage func(transactNexus nexus.Nexus, preimage string, requestReferenceUuid string)

// OnProviderInfo handles a provider info request
type OnProviderInfo func(consumerTransactNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage)

// ConsumerTrackNexus is used for interacting with a consumer via the ConsumerLayer
type ConsumerTrackNexus struct {
	nexus.Nexus
	// invoice event handler
	onInvoice OnInvoice
	// preimage event handler
	onPreimage OnPreimage
	// aclled on provder info
	onProviderInfo OnProviderInfo
}

// ConsumerTrackNexusName is the name of the ConsumerTrackNexusName
const ConsumerTrackNexusName = "ConsumerTrackNexus"

// NewConsumerTransactNexus creates a ConsumerTrackNexus
func NewConsumerTransactNexus(belowNexus nexus.Nexus) *ConsumerTrackNexus {
	c := ConsumerTrackNexus{
		base.NewBaseNexusBelow(ConsumerTrackNexusName, belowNexus),
		nil,
		nil,
		nil,
	}
	belowNexus.SetOnBinMessage(c.OnBinMessage)
	belowNexus.SetOnMessage(c.OnMessage)
	return &c
}

// HandleLayerNotification handles Opinion specific notifications
func (c ConsumerTrackNexus) HandleLayerNotification(msg notification.MoneysocketNotification) {
	if msg.RequestType() == moneysocket_message.NotifyOpinionInvoice {
		notifyMsg := msg.(notification.NotifyInvoice)
		if c.onInvoice != nil {
			c.onInvoice(c, notifyMsg.Bolt11, msg.RequestReferenceUUID())
		}
	} else if msg.RequestType() == moneysocket_message.NotifyPreimage {
		notifyMsg := msg.(notification.NotifyPreimage)
		if c.onPreimage != nil {
			c.onPreimage(c, notifyMsg.Preimage, notifyMsg.RequestReferenceUUID())
		}
	}
}

// IsLayerMessage determines if a message needs to be handled by this layer
func (c ConsumerTrackNexus) IsLayerMessage(msg moneysocket_message.MoneysocketMessage) bool {
	if msg.MessageClass() != moneysocket_message.Notification {
		return false
	}
	notifyMsg := msg.(notification.MoneysocketNotification)
	return notifyMsg.RequestType() == moneysocket_message.NotifyInvoiceNotification ||
		notifyMsg.RequestType() == moneysocket_message.NotifyPreimage
}

// OnMessage handles the message if relevant to this alyer
func (c ConsumerTrackNexus) OnMessage(belowNexus nexus.Nexus, message moneysocket_message.MoneysocketMessage) {
	if !c.IsLayerMessage(message) {
		c.Nexus.OnMessage(belowNexus, message)
	}
}

// OnBinMessage is a callback for a binary message
func (c ConsumerTrackNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	// DO nothing
}

// OnInvoice calls on invoice function
func (c *ConsumerTrackNexus) OnInvoice(transactNexus nexus.Nexus, invoice string, requestReferenceUUID string) {
	if c.onInvoice != nil {
		c.onInvoice(transactNexus, invoice, requestReferenceUUID)
	}
}

// SetOnInvoice sets a function to be called when on invoice is called
func (c *ConsumerTrackNexus) SetOnInvoice(invoice OnInvoice) {
	c.onInvoice = invoice
}

// OnPreImage calls on preimage function
func (c *ConsumerTrackNexus) OnPreImage(transactNexus nexus.Nexus, preimage string, requestReferenceUUID string) {
	if c.onPreimage != nil {
		c.onPreimage(transactNexus, preimage, requestReferenceUUID)
	}
}

// SetOnPreimage sets function to be called when onPreImage is called
func (c *ConsumerTrackNexus) SetOnPreimage(preimage OnPreimage) {
	c.onPreimage = preimage
}

// SetOnProviderInfo sets a function to be called when OnProviderInfo is called
func (c *ConsumerTrackNexus) SetOnProviderInfo(info OnProviderInfo) {
	c.onProviderInfo = info
}

// OnProviderInfo calls onProviderInfo callback
func (c *ConsumerTrackNexus) OnProviderInfo(consumerTransactNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
	if c.onProviderInfo != nil {
		c.onProviderInfo(consumerTransactNexus, msg)
	}
}

// RequestInvoice requests an invoice from the ConsumerTrackNexus
func (c ConsumerTrackNexus) RequestInvoice(msats int64, description string) (uuid string) {
	ri := request.NewRequestInvoice(msats)
	_ = c.Send(ri)
	return ri.UUID()
}

// RequestPay requests a invoice from a callback
func (c ConsumerTrackNexus) RequestPay(bolt11 string) (uuid string) {
	rp := request.NewRequestPay(bolt11)
	_ = c.Send(rp)
	return rp.UUID()
}

var _ nexus.Nexus = &ConsumerTrackNexus{}
