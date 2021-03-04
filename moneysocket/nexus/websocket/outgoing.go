package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/ws/ws_client"
)

const OutgoingSocketName = "OutgoingSocket"

type OutgoingSocket struct {
	ws_client.WebsocketClientProtocol
	nexusHelper.Nexus
	wasAnnounced bool
	// protocol layer coorespond to the socket interface
	FactoryMsProtocolLayer layer.Layer
	// add an outgoing shared seed modules
	OutgoingSharedSeed *beacon.SharedSeed
}

// create a new incoming websocket nexus (accepts request)
func NewOutgoingSocket() OutgoingSocket {
	return OutgoingSocket{
		WebsocketClientProtocol: ws_client.NewBaseWebsocketClient(),
		Nexus:                   base.NewBaseNexus(OutgoingSocketName),
		wasAnnounced:            false,
		OutgoingSharedSeed:      nil,
	}
}

func (i *OutgoingSocket) SharedSeed() *beacon.SharedSeed {
	return i.OutgoingSharedSeed
}

func (i *OutgoingSocket) Send(msg moneysocket_message.MoneysocketMessage) error {
	log.Infof("encoding msg %s", msg)
	ss := i.OutgoingSharedSeed
	msgBytes, err := message.WireEncode(msg, ss)
	if err != nil {
		return err
	}
	return i.WebsocketClientProtocol.SendBin(msgBytes)
}

func (i *OutgoingSocket) SendBin(msg []byte) error {
	return i.WebsocketClientProtocol.SendBin(msg)
}

func (i OutgoingSocket) OnConnecting() {
	log.Info("Websocket connecting")
}

func (i OutgoingSocket) OnConnect(conn *websocket.Conn, r *http.Response) {
	i.WebsocketClientProtocol.OnConnect(conn, r)
	log.Info("Client connecting")
}

func (i *OutgoingSocket) OnOpen() {
	log.Info("websocket connection open")
	i.FactoryMsProtocolLayer.AnnounceNexus(i)
	i.wasAnnounced = true
}

// cooresponds to the nexus interface, handles a message
func (i *OutgoingSocket) OnMessage(belowNexus nexusHelper.Nexus, msg moneysocket_message.MoneysocketMessage) {
	log.Info("websocket nexus got message")
	i.Nexus.OnMessage(belowNexus, msg)
}

// cooresponds to the nexus interface, handles a binary message
func (i *OutgoingSocket) OnBinMessage(belowNexus nexusHelper.Nexus, msg []byte) {
	panic("not yet implemented")
}

func (i *OutgoingSocket) OnWsMessage(payload []byte, isBinary bool) {
	log.Info("outgoing message")
	if isBinary {
		log.Infof("binary payload of %d bytes", len(payload))
		sharedSeed := i.SharedSeed()

		msg, _, err := message.WireDecode(payload, sharedSeed)
		if err != nil {
			log.Infof("could not decode %s", err)
		}
		log.Infof("recv msg: %s", msg)
		i.OnMessage(i, msg)
	} else {
		log.Infof("text payload %s", payload)
		log.Error("text payload is unexpected, dropping")
	}
}

func (i *OutgoingSocket) OnClose(wasClean bool, code int, reason string) {
	log.Infof("websocket connection closed: %s", reason)
	if i.wasAnnounced {
		i.FactoryMsProtocolLayer.RevokeNexus(i)
	}
	i.wasAnnounced = false
}

// assert type is valid socket
var _ ws_client.WebsocketClientProtocol = &OutgoingSocket{}

// assert type is valid nexus
var _ nexusHelper.Nexus = &OutgoingSocket{}
