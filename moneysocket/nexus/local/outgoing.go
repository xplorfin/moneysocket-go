package local

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	base_moneysocket "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// OutgoingLocalNexus is an outgoing local nexus
type OutgoingLocalNexus struct {
	*base.NexusBase
	belowNexus *JoinedLocalNexus
	sharedSeed beacon.SharedSeed
}

// OutgoingLocalNexusName is the name of an outgoing local enxus
const OutgoingLocalNexusName = "OutgoingLocalNexus"

// NewOutgoingLocalNexus creates an outgoing local nexus
func NewOutgoingLocalNexus(belowNexus *JoinedLocalNexus, layer layer.Base, sharedSeed beacon.SharedSeed) *OutgoingLocalNexus {
	bnf := base.NewBaseNexusFull(OutgoingLocalNexusName, belowNexus, layer)
	og := OutgoingLocalNexus{
		NexusBase:  &bnf,
		belowNexus: belowNexus,
	}
	// this needs to be done everywhere
	belowNexus.SetOnBinMessage(og.OnBinMessage)
	belowNexus.SetOnMessage(og.OnMessage)

	og.belowNexus.SetOutgoingNexus(&og)
	og.sharedSeed = sharedSeed
	return &og
}

// OnMessage processes messages for this layer
func (o *OutgoingLocalNexus) OnMessage(belowNexus nexus.Nexus, msg base_moneysocket.MoneysocketMessage) {
	log.Printf("outgoing local nexus got msg %s", msg)
	o.NexusBase.OnMessage(belowNexus, msg)
}

// OnBinMessage processes a mmessage after it's sent
func (o OutgoingLocalNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	log.Printf("outgoing local nexus got raw msg: %d", len(msg))

	proccessedMessage, _, err := message.WireDecode(msg, &o.sharedSeed)
	if err != nil {
		log.Printf("could not decode msg: %s", err)
		o.NexusBase.OnBinMessage(belowNexus, msg)
		return
	}
	o.NexusBase.OnMessage(belowNexus, proccessedMessage)

}

// Send gets the message from a nexus
func (o OutgoingLocalNexus) Send(msg base_moneysocket.MoneysocketMessage) error {
	isEncrypted, msgOrBytes := message.LocalEncode(msg, o.SharedSeed())
	if isEncrypted {
		log.Printf("sending encrypyted: %s", msgOrBytes)
		_ = o.SendBin(msgOrBytes)
	} else {
		o.belowNexus.SendFromOutgoing(msg)
	}
	return nil
}

// SharedSeed gets the shared seed
func (o OutgoingLocalNexus) SharedSeed() *beacon.SharedSeed {
	return &o.sharedSeed
}

// RevokeFromLayer revokes a message from the layer
func (o *OutgoingLocalNexus) RevokeFromLayer() {
	o.Layer.RevokeNexus(o)
}

// TODO we've got a lot more stuff to implement here

var _ nexus.Nexus = &OutgoingLocalNexus{}
var _ compat.RevokableNexus = &OutgoingLocalNexus{}
