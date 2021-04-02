package local

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// IncomingLocalNexus is the local nexus
type IncomingLocalNexus struct {
	*base.NexusBase
}

// IncomingLocalNexusName handles nexuses
const IncomingLocalNexusName = "IncomingLocalNexus"

// NewIncomingLocalNexus creates an incoming local nexus
func NewIncomingLocalNexus(belowNexus *JoinedLocalNexus, layer layer.Base) *IncomingLocalNexus {
	baseNexus := base.NewBaseNexusFull(IncomingLocalNexusName, belowNexus, layer)
	og := IncomingLocalNexus{
		NexusBase: &baseNexus,
	}
	belowNexus.SetOnBinMessage(og.OnBinMessage)
	belowNexus.SetOnMessage(og.OnMessage)
	belowNexus.SetIncomingNexus(&og)
	return &og
}

// OnMessage processes a message
func (i *IncomingLocalNexus) OnMessage(belowNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
	log.Println("incoming local nexus got msg")
	i.NexusBase.OnMessage(belowNexus, msg)
}

// OnBinMessage processes a binary message
func (i *IncomingLocalNexus) OnBinMessage(belowNexus nexus.Nexus, msgBytes []byte) {
	log.Println("incoming local nexus got raw msg")
	i.NexusBase.OnBinMessage(belowNexus, msgBytes)
}

// Send sends a message
func (i *IncomingLocalNexus) Send(msg moneysocket_message.MoneysocketMessage) error {
	belowNexus := (*i.BelowNexus).(*JoinedLocalNexus)

	belowNexus.SendFromIncoming(msg)
	return nil
}

// RevokeFromLayer revokes a layer from a message
func (i *IncomingLocalNexus) RevokeFromLayer() {
	i.Layer.RevokeNexus(i)
}

// TODO we've got a lot more stuff to implement here

var _ nexus.Nexus = &IncomingLocalNexus{}
var _ compat.RevokableNexus = &IncomingLocalNexus{}
