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
	*base.NexusBase
	rendezvousFinishedCb func(nexus.Nexus)
	requestReferenceUUID string
	rendezvousID         string
	directory            *Directory
	// TODO directory
}

const IncomingRendezvousNexusName = "IncomingRendezvousNexus"

func NewIncomingRendezvousNexus(belowNexus nexus.Nexus, layer layer.Layer, directory *Directory) *IncomingRendezvousNexus {
	baseNexus := base.NewBaseNexusFull(IncomingRendezvousNexusName, belowNexus, layer)
	og := IncomingRendezvousNexus{
		NexusBase: &baseNexus,
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
		i.NexusBase.OnMessage(belowNexus, msg)
		return
	}

	req := msg.(request.Rendezvous)
	i.rendezvousID = req.RendezvousID
	i.requestReferenceUUID = req.UUID()

	if i.directory.IsRidPeered(i.rendezvousID) {
		i.InitiateClose()
	}

	i.directory.AddNexus(i, i.rendezvousID)
	peer := i.directory.GetPeerNexus(i.UUID())
	if peer != nil {
		_ = i.Send(notification.NewNotifyRendezvous(i.rendezvousID, i.requestReferenceUUID))
		i.rendezvousFinishedCb(i)
		(*peer).(*IncomingRendezvousNexus).RendezvousAcheived()
	} else {
		_ = i.Send(notification.NewRendezvousNotReady(i.rendezvousID, i.requestReferenceUUID))
	}
}

func (i *IncomingRendezvousNexus) OnBinMessage(belowNexus nexus.Nexus, msgByte []byte) {
	log.Println("rdv nexus got raw message")
	i.NexusBase.OnBinMessage(belowNexus, msgByte)
}

// called by other peer
func (i *IncomingRendezvousNexus) RendezvousAcheived() {
	if !i.directory.IsRidPeered(i.rendezvousID) {
		panic("expected rendezvous to be peered")
	}
	_ = i.Send(notification.NewNotifyRendezvous(i.rendezvousID, i.requestReferenceUUID))
	i.rendezvousFinishedCb(i)
}

func (i *IncomingRendezvousNexus) EndRendezvous() {
	i.directory.RemoveNexus(i)
	_ = i.Send(notification.NewRendezvousEnd(i.rendezvousID, ""))
}

var _ nexus.Nexus = &IncomingRendezvousNexus{}
