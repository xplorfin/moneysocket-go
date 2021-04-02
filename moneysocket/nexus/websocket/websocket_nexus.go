package websocket

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	base2 "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// WebsocketNexusName is a websocket nexus.
const WebsocketNexusName = "WebsocketNexus"

// Nexus is a nexus.
type Nexus struct {
	*base.NexusBase
}

// NewWebsocketNexus creates a websocket nexus.
func NewWebsocketNexus(belowNexus nexus.Nexus, layer layer.Base) *Nexus {
	bnf := base.NewBaseNexusFull(WebsocketNexusName, belowNexus, layer)
	n := Nexus{&bnf}
	n.BelowNexus = &belowNexus
	belowNexus.SetOnMessage(n.OnMessage)
	belowNexus.SetOnBinMessage(n.OnBinMessage)
	// TODO register above nexus here (should really be done all over the place)
	return &n
}

// OnMessage processes a nexus message.
func (o *Nexus) OnMessage(belowNexus nexus.Nexus, msg base2.MoneysocketMessage) {
	log.Println("websocket nexus got msg")
	o.NexusBase.OnMessage(belowNexus, msg)
}

// OnBinMessage processes a binary message.
func (o *Nexus) OnBinMessage(belowNexus nexus.Nexus, msgByte []byte) {
	log.Println("websocket nexus got raw msg")
	o.NexusBase.OnBinMessage(belowNexus, msgByte)
}

var _ nexus.Nexus = &Nexus{}
