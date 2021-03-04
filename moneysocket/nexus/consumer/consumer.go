package consumer

import (
	"log"
	"time"

	"github.com/buger/jsonparser"

	"github.com/xplorfin/moneysocket-go/moneysocket/message"
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
	pingStartTime     time.Time
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

func NewConsumerNexus() ConsumerNexus {
	return ConsumerNexus{
		BaseNexus:    base.NewBaseNexus(ConsumerNexusName),
		donePinging:  make(chan bool, 1),
		pingInterval: time.Second * 3,
	}
}

func (c ConsumerNexus) isLayerMessage(msg []byte) bool {
	msgClass, _ := jsonparser.GetString(msg, message.MessageClass)
	if msgClass != base2.Notification.ToString() {
		return false
	}
	notificationName, _ := jsonparser.GetString(msg, message.NotificationName)
	return containsString([]string{message.NotifyProvider, message.NotifyProviderNotReady, message.NotifyPing}, notificationName)
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
	panic("todo")
	//if !c.isLayerMessage(msg) {
	//	c.BaseNexus.OnMessage(belowNexus, msg)
	//}
	//
	//notificationName := json.GetString(msg, message.NotificationName)
	//switch notificationName {
	//case message.NotifyProvider:
	//	if !c.handshakeFinished {
	//		c.handshakeFinished = true
	//	}
	//	c.ConsumerFinishedCb()
	//case message.NotifyProviderNotReady:
	//	logger.Info("provider not ready, waiting")
	//case message.NotifyPong:
	//	if c.pingStartTime.IsZero() {
	//		return
	//	}
	//	msecs := time.Since(c.pingStartTime).Milliseconds()
	//	if true { //c.OnPing(){ // if onping function is supplied
	//		c.OnPing(c, int(msecs))
	//	}
	//	c.pingStartTime = time.Time{}
	//}
}
func (c *ConsumerNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	c.BaseNexus.OnBinMessage(belowNexus, msg)
}

func (c *ConsumerNexus) StartHandshake(cb ConsumerFinishedCb) {
	c.consumerFinishedCb = cb
	c.Send(request.NewRequestProvider())
}

// send a ping up the chain
func (c *ConsumerNexus) SendPing() {
	c.pingStartTime = time.Now()
	c.Send(request.NewPingRequest())
}

// start pining on a set interval
func (c *ConsumerNexus) StartPinging() {
	tick := time.Tick(3 * time.Second)
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
		case <-tick:
			c.SendPing()
		}
	}
}

func (c *ConsumerNexus) StopPinging() {
	if c.isPinging {
		<-c.donePinging
	}
}

func containsString(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

var _ nexus.Nexus = &ConsumerNexus{}
