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

type OnPingFn func(nexus nexus.Nexus, msecs int)

type ConsumerFinishedCb func(consumerNexus ConsumerNexus)

type ConsumerNexus struct {
	*base.BaseNexus
	handshakeFinished bool
	pingStartTime     *time.Time
	// how often to ping
	pingInterval       time.Duration
	onPing             OnPingFn
	consumerFinishedCb ConsumerFinishedCb
	// wether or not ping loop is currently engaged
	isPinging bool
	// channel for ending ping loop
	donePinging chan bool
}

const ConsumerNexusName = "ConsumerNexus"

func NewConsumerNexus(belowNexus nexus.Nexus) *ConsumerNexus {
	consumerNexus := ConsumerNexus{
		BaseNexus:    base.NewBaseNexus(ConsumerNexusName),
		donePinging:  make(chan bool, 1),
		pingInterval: time.Second * 3,
	}
	belowNexus.SetOnBinMessage(consumerNexus.OnBinMessage)
	belowNexus.SetOnMessage(consumerNexus.OnMessage)
	return &consumerNexus
}

func (c *ConsumerNexus) IsLayerMessage(message base2.MoneysocketMessage) bool {
	if message.MessageClass() != base2.Notification {
		return false
	}
	notif := message.(notification.MoneysocketNotification)
	return notif.MessageClass() == base2.NotifyProvider || notif.MessageClass() == base2.NotifyProviderNotReady || notif.MessageClass() == base2.NotifyPong
}

func (c *ConsumerNexus) ConsumerFinishedCb() {
	// TODO
}

func (c *ConsumerNexus) SetOnPing(fn OnPingFn) {
	c.onPing = fn
}

func (c *ConsumerNexus) OnPing(consumerNexus nexus.Nexus, milliseconds int) {
	if c.onPing != nil {
		c.onPing(consumerNexus, milliseconds)
	}
}

func (c *ConsumerNexus) OnMessage(belowNexus nexus.Nexus, msg base2.MoneysocketMessage) {
	log.Print("consumer nexus got msg")
	if !c.IsLayerMessage(msg) {
		c.BaseNexus.OnMessage(belowNexus, msg)
	}

	notif := msg.(notification.MoneysocketNotification)
	// TODO: switch case?
	if notif.MessageClass() != base2.NotifyProvider {
		if !c.handshakeFinished {
			c.handshakeFinished = true
			c.consumerFinishedCb(*c)
		}
		c.BaseNexus.OnMessage(belowNexus, msg)
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
func (c *ConsumerNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	c.BaseNexus.OnBinMessage(belowNexus, msg)
}

func (c *ConsumerNexus) StartHandshake(cb ConsumerFinishedCb) {
	c.consumerFinishedCb = cb
	_ = c.Send(request.NewRequestProvider())
}

// send a ping up the chain
func (c *ConsumerNexus) SendPing() {
	currentTime := time.Now()
	c.pingStartTime = &currentTime
	_ = c.Send(request.NewPingRequest())
}

// start pining on a set interval
func (c *ConsumerNexus) StartPinging() {
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
			break
		// Got a tick, we should check on doSomething()
		case <-ticker.C:
			c.SendPing()
		}
	}
}

func (c *ConsumerNexus) StopPinging() {
	if c.isPinging {
		<-c.donePinging
	}
}

var _ nexus.Nexus = &ConsumerNexus{}
