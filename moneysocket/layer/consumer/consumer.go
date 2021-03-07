package consumer

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/consumer"
)

type ConsumerLayer struct {
	layer.BaseLayer
	onPing        consumer.OnPingFn
	consumerNexus *consumer.ConsumerNexus
}

func NewConsumerLayer() *ConsumerLayer {
	return &ConsumerLayer{
		layer.NewBaseLayer(),
		nil,
		&consumer.ConsumerNexus{}, // gets overwritten
	}
}

// announce the nexus and start the handshake
func (c *ConsumerLayer) AnnounceNexus(belowNexus nexus.Nexus) {
	c.SetupConsumerNexus(belowNexus)
	c.TrackNexus(c.consumerNexus, belowNexus)
	c.consumerNexus.StartHandshake(c.ConsumerFinishedCb)
}

func (c *ConsumerLayer) SetOnPing(fn consumer.OnPingFn) {
	c.onPing = fn
}

// initialize consumer nexus and tie the onping event back to this layer
func (c *ConsumerLayer) SetupConsumerNexus(belowNexus nexus.Nexus) {
	c.consumerNexus = consumer.NewConsumerNexus(belowNexus)
	c.consumerNexus.SetOnPing(c.OnPing)
}

func (c *ConsumerLayer) RegisterAboveLayer(belowLayer layer.Layer) {
	belowLayer.SetOnAnnounce(c.AnnounceNexus)
	belowLayer.SetOnRevoke(c.RevokeNexus)
}

func (c *ConsumerLayer) RevokeNexus(belowNexus nexus.Nexus) {
	// TODO add error handling
	belowUuid, _ := c.NexusByBelow.Get(belowNexus.Uuid())
	consumerNexus, _ := c.Nexuses.Get(belowUuid)
	c.BaseLayer.RevokeNexus(consumerNexus)
	castedNexus := consumerNexus.(*consumer.ConsumerNexus)
	castedNexus.StopPinging()

}

// event fired on ping
func (c *ConsumerLayer) OnPing(consumerNexus nexus.Nexus, milliseconds int) {
	if c.onPing != nil {
		c.onPing(consumerNexus, milliseconds)
	}
}

// consume finished
func (c *ConsumerLayer) ConsumerFinishedCb(consumerNexus consumer.ConsumerNexus) {
	c.TrackNexusAnnounced(&consumerNexus)
	c.SendLayerEvent(&consumerNexus, message.NexusAnnounced)
	if c.OnAnnounce != nil {
		c.OnAnnounce(&consumerNexus)
	}
	consumerNexus.StartPinging()
}

var _ layer.Layer = &ConsumerLayer{}
