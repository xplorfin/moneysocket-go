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

type FinishedCb func(consumerNexus Nexus)

type Nexus struct {
	*base.NexusBase
	handshakeFinished bool
	pingStartTime     *time.Time
	// how often to ping
	pingInterval       time.Duration
	onPing             OnPingFn
	consumerFinishedCb FinishedCb
	// wether or not ping loop is currently engaged
	isPinging bool
	// channel for ending ping loop
	donePinging chan bool
}

const NexusName = "ConsumerNexus"

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

func (c *Nexus) IsLayerMessage(message base2.MoneysocketMessage) bool {
	if message.MessageClass() != base2.Notification {
		return false
	}
	notif := message.(notification.MoneysocketNotification)
	return notif.MessageClass() == base2.NotifyProvider || notif.MessageClass() == base2.NotifyProviderNotReady || notif.MessageClass() == base2.NotifyPong
}

func (c *Nexus) ConsumerFinishedCb() {
	// TODO
}

func (c *Nexus) SetOnPing(fn OnPingFn) {
	c.onPing = fn
}

func (c *Nexus) OnPing(consumerNexus nexus.Nexus, milliseconds int) {
	if c.onPing != nil {
		c.onPing(consumerNexus, milliseconds)
	}
}

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
func (c *Nexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	c.NexusBase.OnBinMessage(belowNexus, msg)
}

func (c *Nexus) StartHandshake(cb FinishedCb) {
	c.consumerFinishedCb = cb
	_ = c.Send(request.NewRequestProvider())
}

// send a ping up the chain
func (c *Nexus) SendPing() {
	currentTime := time.Now()
	c.pingStartTime = &currentTime
	_ = c.Send(request.NewPingRequest())
}

// start pining on a set interval
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

func (c *Nexus) StopPinging() {
	if c.isPinging {
		<-c.donePinging
	}
}

var _ nexus.Nexus = &Nexus{}
