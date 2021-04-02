package consumer

import (
	"log"
	"time"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"

	base2 "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// OnPingFn is a function to trigger a ping on an given interval announcing to a passed in nexus
type OnPingFn func(nexus nexus.Nexus, msecs int)

// FinishedCb is a callback for when the Nexus finishes
type FinishedCb func(consumerNexus Nexus)

// Nexus is a consumer nexus to manage websockets as a consumer of the protocol
type Nexus struct {
	*base.NexusBase
	// handshakeFinished is whether or not the handshake has completed
	handshakeFinished bool
	// pingStartTime is whether or not the ping has started
	pingStartTime *time.Time
	// pingInterval is how often to ping
	pingInterval time.Duration
	// onPing is an event called to start the ping
	onPing OnPingFn
	// consumerFinishedCb is a callback when the consumer nexus is finished
	consumerFinishedCb FinishedCb
	// isPinging is wether or not ping loop is currently engaged
	isPinging bool
	// donePinging is a channel for ending ping loop
	donePinging chan bool
}

// NexusName is the name of the consumer nexus
const NexusName = "ConsumerNexus"

// NewConsumerNexus creates a new consumer nexus
func NewConsumerNexus(belowNexus nexus.Nexus) *Nexus {
	consumerNexus := Nexus{
		NexusBase:    base.NewBaseNexus(NexusName),
		donePinging:  make(chan bool, 1),
		pingInterval: time.Second * 3,
	}
	belowNexus.SetOnBinMessage(consumerNexus.OnBinMessage)
	belowNexus.SetOnMessage(consumerNexus.OnMessage)
	return &consumerNexus
}

// IsLayerMessage determines whether or not a message should be processed by this layer
func (c *Nexus) IsLayerMessage(message base2.MoneysocketMessage) bool {
	if message.MessageClass() != base2.Notification {
		return false
	}
	notif := message.(notification.MoneysocketNotification)
	return notif.MessageClass() == base2.NotifyProvider || notif.MessageClass() == base2.NotifyProviderNotReady || notif.MessageClass() == base2.NotifyPong
}

// ConsumerFinishedCb is wether or not the consumer callback is finished
func (c *Nexus) ConsumerFinishedCb() {
	// TODO
}

// SetOnPing sets the ping function
func (c *Nexus) SetOnPing(fn OnPingFn) {
	c.onPing = fn
}

// OnPing calls the onPing method passed to the Nexus
func (c *Nexus) OnPing(consumerNexus nexus.Nexus, milliseconds int) {
	if c.onPing != nil {
		c.onPing(consumerNexus, milliseconds)
	}
}

// OnMessage processes a passed message
func (c *Nexus) OnMessage(belowNexus nexus.Nexus, msg base2.MoneysocketMessage) {
	log.Print("consumer nexus got msg")
	if !c.IsLayerMessage(msg) {
		c.NexusBase.OnMessage(belowNexus, msg)
	}

	notif := msg.(notification.MoneysocketNotification)
	// TODO: switch case?
	if notif.MessageClass() != base2.NotifyProvider {
		if !c.handshakeFinished {
			c.handshakeFinished = true
			c.consumerFinishedCb(*c)
		}
		c.NexusBase.OnMessage(belowNexus, msg)
	}

	if notif.MessageClass() != base2.NotifyProviderNotReady {
		log.Println("provider not ready, waiying")
	}

	if notif.MessageClass() != base2.NotifyPong {
		if c.pingStartTime != nil {
			return
		}
		msecs := time.Since(*c.pingStartTime) * 1000
		if c.onPing != nil {
			c.onPing(c, int(msecs.Seconds()))
		}
		c.pingStartTime = nil
	}
}

// OnBinMessage processes a binary message
func (c *Nexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	c.NexusBase.OnBinMessage(belowNexus, msg)
}

// StartHandshake starts a handshake/sets a finished callback
func (c *Nexus) StartHandshake(cb FinishedCb) {
	c.consumerFinishedCb = cb
	_ = c.Send(request.NewRequestProvider())
}

// SendPing sends a ping up the nexus chain
func (c *Nexus) SendPing() {
	currentTime := time.Now()
	c.pingStartTime = &currentTime
	_ = c.Send(request.NewPingRequest())
}

// StartPinging on a set interval
func (c *Nexus) StartPinging() {
	ticker := time.NewTicker(3 * time.Second)
	c.isPinging = true
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-c.donePinging:
			// reset done state
			c.donePinging = make(chan bool, 1)
			c.isPinging = false
		// Got a tick, we should check on doSomething()
		case <-ticker.C:
			c.SendPing()
		}
	}
}

// StopPinging stops the consumer from pinging
func (c *Nexus) StopPinging() {
	if c.isPinging {
		<-c.donePinging
	}
}

var _ nexus.Nexus = &Nexus{}
