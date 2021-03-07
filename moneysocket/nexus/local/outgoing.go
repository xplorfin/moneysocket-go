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

type OutgoingLocalNexus struct {
	*base.BaseNexus
	belowNexus *JoinedLocalNexus
	sharedSeed beacon.SharedSeed
}

const OutgoingLocalNexusName = "OutgoingLocalNexus"

func NewOutgoingLocalNexus(belowNexus *JoinedLocalNexus, layer layer.Layer, sharedSeed beacon.SharedSeed) *OutgoingLocalNexus {
	bnf := base.NewBaseNexusFull(OutgoingLocalNexusName, belowNexus, layer)
	og := OutgoingLocalNexus{
		BaseNexus:  &bnf,
		belowNexus: belowNexus,
	}
	// this needs to be done everywhere
	belowNexus.SetOnBinMessage(og.OnBinMessage)
	belowNexus.SetOnMessage(og.OnMessage)

	og.belowNexus.SetOutgoingNexus(&og)
	og.sharedSeed = sharedSeed
	return &og
}

func (o *OutgoingLocalNexus) OnMessage(belowNexus nexus.Nexus, msg base_moneysocket.MoneysocketMessage) {
	log.Printf("outgoing local nexus got msg %s", msg)
	o.BaseNexus.OnMessage(belowNexus, msg)
}

func (o OutgoingLocalNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	log.Printf("outgoing local nexus got raw msg: %d", len(msg))

	proccessedMessage, _, err := message.WireDecode(msg, &o.sharedSeed)
	if err != nil {
		log.Printf("could not decode msg: %s", err)
		o.BaseNexus.OnBinMessage(belowNexus, msg)
		return
	}
	o.BaseNexus.OnMessage(belowNexus, proccessedMessage)

}

func (o OutgoingLocalNexus) Send(msg base_moneysocket.MoneysocketMessage) error {
	isEncrypted, msgOrBytes := message.LocalEncode(msg, o.SharedSeed())
	if isEncrypted {
		log.Printf("sending encrypyted: %s", msgOrBytes)
		o.SendBin(msgOrBytes)
	} else {
		o.belowNexus.SendFromOutgoing(msg)
	}
	return nil
}

func (o OutgoingLocalNexus) SharedSeed() *beacon.SharedSeed {
	return &o.sharedSeed
}

func (o *OutgoingLocalNexus) RevokeFromLayer() {
	o.Layer.RevokeNexus(o)
}

// TODO we've got a lot more stuff to implement here

var _ nexus.Nexus = &OutgoingLocalNexus{}
var _ compat.RevokableNexus = &OutgoingLocalNexus{}
