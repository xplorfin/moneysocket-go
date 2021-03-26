package rendezvous

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	message_base "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	notification2 "github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

type OutgoingRendezvousNexus struct {
	*base.NexusBase
	rendezvousFinishedCb func(nexus.Nexus)
}

const OutgoingRendezvousNexusName = "OutgoingRendezvousNexus"

func NewOutgoingRendezvousNexus(belowNexus nexus.Nexus, layer layer.LayerBase) *OutgoingRendezvousNexus {
	bnf := base.NewBaseNexusFull(OutgoingRendezvousNexusName, belowNexus, layer)
	og := OutgoingRendezvousNexus{
		NexusBase: &bnf,
	}
	belowNexus.SetOnBinMessage(og.OnBinMessage)
	belowNexus.SetOnMessage(og.OnMessage)

	return &og
}

func (o OutgoingRendezvousNexus) IsLayerMessage(msg message_base.MoneysocketMessage) bool {
	if msg.MessageClass() == message_base.Notification {
		return false
	}
	notification := msg.(notification2.MoneysocketNotification)
	return notification.RequestType() == message_base.NotifyRendezvous ||
		notification.RequestType() == message_base.NotifyRendezvousNotReadyNotification ||
		notification.RequestType() == message_base.NotifyRendezvousEndNotification
}

func (o *OutgoingRendezvousNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	log.Println("rdv nexus got raw message") //apparently this shouldn't happen
}

func (o *OutgoingRendezvousNexus) OnMessage(belowNexus nexus.Nexus, msg message_base.MoneysocketMessage) {
	log.Printf("outgoing rdv nexus got msg %s", msg)
	if !o.IsLayerMessage(msg) {
		o.NexusBase.OnMessage(belowNexus, msg)
	}

	notif := msg.(notification2.MoneysocketNotification)
	switch notif.RequestType() {
	case message_base.NotifyRendezvous:
		log.Println("rendezvous ready, notifying")
		o.rendezvousFinishedCb(o)
	case message_base.NotifyRendezvousNotReadyNotification:
		log.Println("rendezvous not ready, waiting")
	case message_base.NotifyRendezvousEndNotification:
		log.Println("rendezvous ended")
		o.InitiateClose()
	}
}

func (o *OutgoingRendezvousNexus) StartRendezvous(rendevousID string, rendezvousFinishedCb func(nexus nexus.Nexus)) {
	o.rendezvousFinishedCb = rendezvousFinishedCb
	rendezvousRequest := request.NewRendezvousRequest(rendevousID)
	_ = o.Send(rendezvousRequest)
}

var _ nexus.Nexus = &OutgoingRendezvousNexus{}
