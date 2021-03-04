package websocket

import (
	"net/http"

	"github.com/prometheus/common/log"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message"
	base2 "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/ws/ws_server"
)

type IncomingSocket struct {
	ws_server.WebSocketServerProtocol
	nexusHelper.Nexus
	wasAnnounced bool
	// protocol layer coorespond to the socket interface
	FactoryMsProtocolLayer layer.Layer
	FactoryMsSharedSeed    *beacon.SharedSeed
}

const IncomingSocketName = "IncomingSocketName"

// create a new incoming websocket nexus (accepts request)
func NewIncomingSocket() IncomingSocket {
	return IncomingSocket{
		WebSocketServerProtocol: ws_server.NewBaseWebsocketService(),
		Nexus:                   base.NewBaseNexus(IncomingSocketName),
		wasAnnounced:            false,
	}
}

func (i IncomingSocket) OnConnecting(r *http.Request) {
	log.Info("Websocket connecting")
}

func (i IncomingSocket) OnConnect(r *http.Request) {
	log.Info("Client connecting")
}

func (i *IncomingSocket) OnOpen() {
	log.Info("websocket connection open")
	i.FactoryMsProtocolLayer.AnnounceNexus(i)
	i.wasAnnounced = true
}

func (i *IncomingSocket) Send(msg base2.MoneysocketMessage) error {
	log.Infof("encoding msg %s", msg.MessageClass().ToString())
	ss := i.SharedSeed()
	msgBytes, err := message.WireEncode(msg, ss)
	if err != nil {
		return err
	}
	return i.SendBin(msgBytes)
}

func (i *IncomingSocket) SendBin(rawMsg []byte) error {
	return i.WebSocketServerProtocol.SendMessage(rawMsg)
}

// cooresponds to the nexus interface, handles a message
func (i *IncomingSocket) OnMessage(belowNexus nexusHelper.Nexus, moneysocketMessage base2.MoneysocketMessage) {
	log.Info("websocket nexus got message")
	i.Nexus.OnMessage(belowNexus, moneysocketMessage)
}

// cooresponds to the nexus interface, handles a binary message
func (i *IncomingSocket) OnBinMessage(belowNexus nexusHelper.Nexus, msg []byte) {
	i.Nexus.OnBinMessage(belowNexus, msg)
}

func (i *IncomingSocket) OnWsMessage(payload []byte, isBinary bool) {
	if isBinary {
		log.Infof("binary payload of %d bytes", len(payload))
		sharedSeed := i.SharedSeed()

		// this needs to be flipped
		if sharedSeed != nil && message.IsCypherText(payload) {
			i.Nexus.OnBinMessage(i, payload)
			return
		}

		msg, _, err := message.WireDecode(payload, sharedSeed)
		if err != nil {
			log.Infof("could not decode %s", err)
		}
		log.Infof("recv msg: %s", msg)
		msg, _, _ = message.MessageFromText(payload)
		i.OnMessage(i, msg)
	} else {
		log.Infof("text payload %s", payload)
		log.Error("text payload is unexpected, dropping")
	}
}

func (i *IncomingSocket) OnClose(wasClean bool, code int, reason string) {
	defer func() {
		if r := recover(); r != nil {
			log.Info("failed to revoke nexus")
		}
	}()
	log.Infof("websocket connection closed: %s", reason)
	if i.wasAnnounced {
		i.FactoryMsProtocolLayer.RevokeNexus(i)
	}
	i.wasAnnounced = false
}

func (i *IncomingSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws_server.ServeHTTP(i, w, r)
}

func (i *IncomingSocket) SharedSeed() *beacon.SharedSeed {
	return i.FactoryMsSharedSeed
}

func (i *IncomingSocket) InitiateClose() {
	i.Cancel()()
}

// assert type is valid socket
var _ ws_server.WebSocketServerProtocol = &IncomingSocket{}

// assert type is valid nexus
var _ nexusHelper.Nexus = &IncomingSocket{}
