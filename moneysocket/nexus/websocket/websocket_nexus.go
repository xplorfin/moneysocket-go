package websocket

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"log"

	base2 "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

const WebsocketNexusName = "WebsocketNexus"

type WebsocketNexus struct {
	*base.BaseNexus
}

func NewWebsocketNexus(belowNexus nexus.Nexus, layer layer.Layer) *WebsocketNexus {
	bnf := base.NewBaseNexusFull(WebsocketNexusName, belowNexus, layer)
	n := WebsocketNexus{&bnf}
	n.BelowNexus = &belowNexus
	belowNexus.SetOnMessage(n.OnMessage)
	belowNexus.SetOnBinMessage(n.OnBinMessage)
	// TODO register above nexus here (should really be done all over the place)
	return &n
}

func (o *WebsocketNexus) OnMessage(belowNexus nexus.Nexus, msg base2.MoneysocketMessage) {
	log.Println("websocket nexus got msg")
	o.BaseNexus.OnMessage(belowNexus, msg)
}

func (o *WebsocketNexus) OnBinMessage(belowNexus nexus.Nexus, msgByte []byte) {
	log.Println("websocket nexus got raw msg")
	o.BaseNexus.OnBinMessage(belowNexus, msgByte)
}

var _ nexus.Nexus = &WebsocketNexus{}
