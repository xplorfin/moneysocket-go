package stack

import (
	"fmt"
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/rendezvous"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer/websocket"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
)

// outgoing consumer stack provides an interoperable interface to connect to a given beacon
// currently it only supports the websocket layer
type OutgoingConsumerStack struct {
	*ConsumerStack
	websocketLayer  *websocket.OutgoingWebsocketLayer
	rendezvousLayer *rendezvous.OutgoingRendezvousLayer
}

// create and initialize an outgoing consumer stack
func NewOutgoingConsumerStack() *OutgoingConsumerStack {
	outgoingConsumerStack := OutgoingConsumerStack{
		NewConsumerStack(),
		websocket.NewOutgoingWebsocketLayer(),
		rendezvous.NewOutgoingRendezvousLayer(),
	}
	// TODO clean this up: don't double initialize
	outgoingConsumerStack.SetupWebsocketLayer()
	outgoingConsumerStack.SetupRendezvousLayer(outgoingConsumerStack.websocketLayer)
	outgoingConsumerStack.SetupConsumerLayer(outgoingConsumerStack.rendezvousLayer)
	outgoingConsumerStack.SetupTransactLayer(outgoingConsumerStack.consumerLayer)
	return &outgoingConsumerStack
}

func (o *OutgoingConsumerStack) SetupRendezvousLayer(belowLayer layer.LayerBase) {
	o.rendezvousLayer = rendezvous.NewOutgoingRendezvousLayer()
	o.rendezvousLayer.RegisterAboveLayer(belowLayer)
	o.rendezvousLayer.RegisterLayerEvent(o.sendStackEvent, message.OutgoingRendezvous)
}

func (o *OutgoingConsumerStack) SetupWebsocketLayer() {
	o.websocketLayer = websocket.NewOutgoingWebsocketLayer()
	o.websocketLayer.RegisterLayerEvent(o.sendStackEvent, message.OutgoingWebsocket)
}

// connect to the given beacon
func (o *OutgoingConsumerStack) DoConnect(connectBeacon beacon.Beacon) error {
	log.Println("stack connect called")
	loc := connectBeacon.Locations()[0]
	sharedSeed := connectBeacon.GetSharedSeed()
	if loc.Type() != util.WebsocketLocationTLVType {
		panic(fmt.Errorf("location type %d not yet implemented", loc.Type()))
	}
	wsLocation := loc.(location.WebsocketLocation)
	_, err := o.websocketLayer.Connect(wsLocation, &sharedSeed)
	return err
}
