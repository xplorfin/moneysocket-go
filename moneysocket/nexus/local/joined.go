package local

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	base2 "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

const JoinedLocalNexusName = "JoinedLocalNexus"

type JoinedLocalNexus struct {
	*base.BaseNexus
	outgoingNexus compat.RevokableNexus
	incomingNexus compat.RevokableNexus

	// TODO event listeners
}

func (j *JoinedLocalNexus) SendFromOutgoing(msg base2.MoneysocketMessage) {
	log.Printf("from outgoing %s", msg)
	j.incomingNexus.OnMessage(j, msg)
}

func (j *JoinedLocalNexus) SendFromIncoming(msg base2.MoneysocketMessage) {
	log.Printf("from incoming %s", msg)
	j.outgoingNexus.OnMessage(j, msg)
}

func (j *JoinedLocalNexus) SendBinFromIncoming(msg []byte) {
	log.Printf("raw from incoming: %d", len(msg))
	j.outgoingNexus.OnBinMessage(j, msg)
}

func NewJoinedLocalNexus() *JoinedLocalNexus {
	bn := base.NewBaseNexus(JoinedLocalNexusName)
	return &JoinedLocalNexus{
		BaseNexus:     bn,
		outgoingNexus: nil,
		incomingNexus: nil,
	}
}

func (n *JoinedLocalNexus) SetIncomingNexus(incomingNexus compat.RevokableNexus) {
	n.incomingNexus = incomingNexus
}

func (n *JoinedLocalNexus) SetOutgoingNexus(outgoingNexus compat.RevokableNexus) {
	n.outgoingNexus = outgoingNexus
}

func (n *JoinedLocalNexus) InitiateClose() {
	n.incomingNexus.InitiateClose()
	n.outgoingNexus.InitiateClose()
}

var _ nexusHelper.Nexus = &JoinedLocalNexus{}
