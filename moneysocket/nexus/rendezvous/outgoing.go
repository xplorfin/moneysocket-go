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

// TODO implement this
type OutgoingRendezvousNexus struct {
	*base.BaseNexus
	rendezvousFinishedCb func(nexus.Nexus)
}

const OutgoingRendezvousNexusName = "OutgoingRendezvousNexus"

func NewOutgoingRendezvousNexus(belowNexus nexus.Nexus, layer layer.Layer) OutgoingRendezvousNexus {
	bnf := base.NewBaseNexusFull(OutgoingRendezvousNexusName, belowNexus, layer)
	og := OutgoingRendezvousNexus{
		BaseNexus: &bnf,
	}
	belowNexus.SetOnBinMessage(og.OnBinMessage)
	belowNexus.SetOnMessage(og.OnMessage)

	return og
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
		o.BaseNexus.OnMessage(belowNexus, msg)
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

func (o *OutgoingRendezvousNexus) StartRendezvous(rendevousId string, rendezvousFinishedCb func(nexus2 nexus.Nexus)) {
	o.rendezvousFinishedCb = rendezvousFinishedCb
	rendezvousRequest := request.NewRendezvousRequest(rendevousId)
	o.Send(rendezvousRequest)
}

var _ nexus.Nexus = &OutgoingRendezvousNexus{}
