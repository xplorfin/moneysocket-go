package rendezvous

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	message_base "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

type IncomingRendezvousNexus struct {
	*base.BaseNexus
	rendezvousFinishedCb func(nexus.Nexus)
	requestReferenceUuid string
	rendezvousId         string
	providerFinishedCb   func(nexus.Nexus)
	directory            *RendezvousDirectory
	// TODO directory
}

const IncomingRendezvousNexusName = "IncomingRendezvousNexus"

func NewIncomingRendezvousNexus(belowNexus nexus.Nexus, layer layer.Layer, directory *RendezvousDirectory) *IncomingRendezvousNexus {
	baseNexus := base.NewBaseNexusFull(IncomingRendezvousNexusName, belowNexus, layer)
	og := IncomingRendezvousNexus{
		BaseNexus: &baseNexus,
		directory: directory,
	}
	belowNexus.SetOnBinMessage(og.OnBinMessage)
	belowNexus.SetOnMessage(og.OnMessage)

	return &og
}

func (i *IncomingRendezvousNexus) IsLayerMessage(msg message_base.MoneysocketMessage) bool {
	if msg.MessageClass() != message_base.Request {
		return false
	}
	req := msg.(request.MoneysocketRequest)
	return req.MessageType() == message_base.RendezvousRequest
}

func (i *IncomingRendezvousNexus) WaitForRendezvous(rendezvousFinishedCb func(nexus.Nexus)) {
	i.rendezvousFinishedCb = rendezvousFinishedCb
}

func (i *IncomingRendezvousNexus) OnMessage(belowNexus nexus.Nexus, msg message_base.MoneysocketMessage) {
	log.Println("rdv nexus got message")
	if !i.IsLayerMessage(msg) {
		i.BaseNexus.OnMessage(belowNexus, msg)
		return
	}

	req := msg.(request.RequestRendezvous)
	i.rendezvousId = req.RendezvousId
	i.requestReferenceUuid = req.Uuid()

	if i.directory.IsRidPeered(i.rendezvousId) {
		i.InitiateClose()
	}

	i.directory.AddNexus(i, i.rendezvousId)
	peer := i.directory.GetPeerNexus(i.UUID())
	if peer != nil {
		_ = i.Send(notification.NewNotifyRendezvous(i.rendezvousId, i.requestReferenceUuid))
		i.rendezvousFinishedCb(i)
		(*peer).(*IncomingRendezvousNexus).RendezvousAcheived()
	} else {
		_ = i.Send(notification.NewRendezvousNotReady(i.rendezvousId, i.requestReferenceUuid))
	}
}

func (i *IncomingRendezvousNexus) OnBinMessage(belowNexus nexus.Nexus, msgByte []byte) {
	log.Println("rdv nexus got raw message")
	i.BaseNexus.OnBinMessage(belowNexus, msgByte)
}

// called by other peer
func (i *IncomingRendezvousNexus) RendezvousAcheived() {
	if !i.directory.IsRidPeered(i.rendezvousId) {
		panic("expected rendezvous to be peered")
	}
	i.Send(notification.NewNotifyRendezvous(i.rendezvousId, i.requestReferenceUuid))
	i.rendezvousFinishedCb(i)
}

func (i *IncomingRendezvousNexus) EndRendezvous() {
	i.directory.RemoveNexus(i)
	i.Send(notification.NewRendezvousEnd(i.rendezvousId, ""))
}

var _ nexus.Nexus = &IncomingRendezvousNexus{}
