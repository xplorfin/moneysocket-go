package local

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

type IncomingLocalNexus struct {
	*base.BaseNexus
}

const IncomingLocalNexusName = "IncomingLocalNexus"

func NewIncomingLocalNexus(belowNexus *JoinedLocalNexus, layer layer.Layer) *IncomingLocalNexus {
	baseNexus := base.NewBaseNexusFull(IncomingLocalNexusName, belowNexus, layer)
	og := IncomingLocalNexus{
		BaseNexus: &baseNexus,
	}
	belowNexus.SetOnBinMessage(og.OnBinMessage)
	belowNexus.SetOnMessage(og.OnMessage)
	belowNexus.SetIncomingNexus(&og)
	return &og
}

func (i *IncomingLocalNexus) OnMessage(belowNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
	log.Println("incoming local nexus got msg")
	i.BaseNexus.OnMessage(belowNexus, msg)
}

func (i *IncomingLocalNexus) OnBinMessage(belowNexus nexus.Nexus, msgBytes []byte) {
	log.Println("incoming local nexus got raw msg")
	i.BaseNexus.OnBinMessage(belowNexus, msgBytes)
}

func (i *IncomingLocalNexus) Send(msg moneysocket_message.MoneysocketMessage) error {
	belowNexus := (*i.BelowNexus).(*JoinedLocalNexus)

	belowNexus.SendFromIncoming(msg)
	return nil
}

func (i *IncomingLocalNexus) RevokeFromLayer() {
	i.Layer.RevokeNexus(i)
}

// TODO we've got a lot more stuff to implement here

var _ nexus.Nexus = &IncomingLocalNexus{}
var _ compat.RevokableNexus = &IncomingLocalNexus{}
