package local

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	base2 "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// JoinedLocalNexusName is a joined local nexus.
const JoinedLocalNexusName = "JoinedLocalNexus"

// JoinedLocalNexus is a joined local nexus.
type JoinedLocalNexus struct {
	*base.NexusBase
	outgoingNexus compat.RevokableNexus
	incomingNexus compat.RevokableNexus

	// TODO event listeners
}

// NewJoinedLocalNexus sets an incoming nexus.
func NewJoinedLocalNexus() *JoinedLocalNexus {
	bn := base.NewBaseNexus(JoinedLocalNexusName)
	return &JoinedLocalNexus{
		NexusBase:     bn,
		outgoingNexus: nil,
		incomingNexus: nil,
	}
}

// SendFromOutgoing sends a message from an outgoing message.
func (j *JoinedLocalNexus) SendFromOutgoing(msg base2.MoneysocketMessage) {
	log.Printf("from outgoing %s", msg)
	j.incomingNexus.OnMessage(j, msg)
}

// SendFromIncoming sends a message from an incoming nexus.
func (j *JoinedLocalNexus) SendFromIncoming(msg base2.MoneysocketMessage) {
	log.Printf("from incoming %s", msg)
	j.outgoingNexus.OnMessage(j, msg)
}

// SendBinFromIncoming sends a binary message.
func (j *JoinedLocalNexus) SendBinFromIncoming(msg []byte) {
	log.Printf("raw from incoming: %d", len(msg))
	j.outgoingNexus.OnBinMessage(j, msg)
}

// SetIncomingNexus sets the incoming nexus.
func (j *JoinedLocalNexus) SetIncomingNexus(incomingNexus compat.RevokableNexus) {
	j.incomingNexus = incomingNexus
}

// SetOutgoingNexus sets the outgoing nexus.
func (j *JoinedLocalNexus) SetOutgoingNexus(outgoingNexus compat.RevokableNexus) {
	j.outgoingNexus = outgoingNexus
}

// InitiateClose initiates a close w/ the incoming/outgoing nexus.
func (j *JoinedLocalNexus) InitiateClose() {
	j.incomingNexus.InitiateClose()
	j.outgoingNexus.InitiateClose()
}

var _ nexusHelper.Nexus = &JoinedLocalNexus{}
