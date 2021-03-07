package websocket

import (
	uuid "github.com/satori/go.uuid"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/ws/ws_client"
)

const OutgoingSocketName = "OutgoingSocket"

type OutgoingSocket struct {
	ws_client.WebsocketClientProtocol
	wasAnnounced bool
	// protocol layer coorespond to the socket interface
	FactoryMsProtocolLayer layer.Layer
	// add an outgoing shared seed modules
	OutgoingSharedSeed *beacon.SharedSeed
	// name of the nexus (stored for debugging)
	name string
	// uuid of the nexus
	uuid uuid.UUID
	// on message
	onMessage nexusHelper.OnMessage
	// on bin message
	onBinMessage nexusHelper.OnBinMessage
}

// create a new incoming websocket nexus (accepts request)
func NewOutgoingSocket() OutgoingSocket {
	return OutgoingSocket{
		WebsocketClientProtocol: ws_client.NewBaseWebsocketClient(),
		wasAnnounced:            false,
		OutgoingSharedSeed:      nil,
		name:                    OutgoingSocketName,
		uuid:                    uuid.NewV4(),
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
	i.onMessage(belowNexus, msg)
}

// cooresponds to the nexus interface, handles a binary message
func (i *OutgoingSocket) OnBinMessage(belowNexus nexusHelper.Nexus, msg []byte) {
	i.onBinMessage(belowNexus, msg)
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

func (i *OutgoingSocket) Uuid() uuid.UUID {
	return i.uuid
}

func (i *OutgoingSocket) IsEqual(n nexusHelper.Nexus) bool {
	panic("implement me")
}

func (i *OutgoingSocket) GetDownwardNexusList() []nexusHelper.Nexus {
	panic("implement me")
}

func (i *OutgoingSocket) InitiateClose() {
	panic("implement me")
}

func (i *OutgoingSocket) Name() string {
	return i.name
}

func (i *OutgoingSocket) SetOnMessage(messageFunc nexusHelper.OnMessage) {
	i.onMessage = messageFunc
}

func (i *OutgoingSocket) SetOnBinMessage(messageBinFunc nexusHelper.OnBinMessage) {
	i.onBinMessage = messageBinFunc
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
