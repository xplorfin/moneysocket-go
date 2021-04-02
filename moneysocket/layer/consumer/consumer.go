package consumer

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/consumer"
)

// Layer is the ConsumerLayer
type Layer struct {
	layer.BaseLayer
	onPing        consumer.OnPingFn
	consumerNexus *consumer.Nexus
}

// NewConsumerLayer creates a new consumer layer
func NewConsumerLayer() *Layer {
	return &Layer{
		layer.NewBaseLayer(),
		nil,
		&consumer.Nexus{}, // gets overwritten
	}
}

// AnnounceNexus creates a new consumer.ConsumerNexus and starts the handshake
func (c *Layer) AnnounceNexus(belowNexus nexus.Nexus) {
	c.SetupConsumerNexus(belowNexus)
	c.TrackNexus(c.consumerNexus, belowNexus)
	c.consumerNexus.StartHandshake(c.ConsumerFinishedCb)
}

// SetOnPing sets up the on ping function
func (c *Layer) SetOnPing(fn consumer.OnPingFn) {
	c.onPing = fn
}

// SetupConsumerNexus initializes consumer nexus and tie the onping event back to this layer
func (c *Layer) SetupConsumerNexus(belowNexus nexus.Nexus) {
	c.consumerNexus = consumer.NewConsumerNexus(belowNexus)
	c.consumerNexus.SetOnPing(c.OnPing)
}

// RegisterAboveLayer registers the current nexuses announce/revoke nexuses to the below layer
func (c *Layer) RegisterAboveLayer(belowLayer layer.Base) {
	belowLayer.SetOnAnnounce(c.AnnounceNexus)
	belowLayer.SetOnRevoke(c.RevokeNexus)
}

// RevokeNexus is a revoked nexus
func (c *Layer) RevokeNexus(belowNexus nexus.Nexus) {
	// TODO add error handling
	belowUUID, _ := c.NexusByBelow.Get(belowNexus.UUID())
	consumerNexus, _ := c.Nexuses.Get(belowUUID)
	c.BaseLayer.RevokeNexus(consumerNexus)
	castedNexus := consumerNexus.(*consumer.Nexus)
	castedNexus.StopPinging()

}

// OnPing is an event fired on ping
func (c *Layer) OnPing(consumerNexus nexus.Nexus, milliseconds int) {
	if c.onPing != nil {
		c.onPing(consumerNexus, milliseconds)
	}
}

// ConsumerFinishedCb is a call back for when a consumer is finished
func (c *Layer) ConsumerFinishedCb(consumerNexus consumer.Nexus) {
	c.TrackNexusAnnounced(&consumerNexus)
	c.SendLayerEvent(&consumerNexus, message.NexusAnnounced)
	if c.OnAnnounce != nil {
		c.OnAnnounce(&consumerNexus)
	}
	consumerNexus.StartPinging()
}

var _ layer.Base = &Layer{}
