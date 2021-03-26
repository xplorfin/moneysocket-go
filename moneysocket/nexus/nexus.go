package nexus

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type OnMessage = func(belowNexus Nexus, msg base.MoneysocketMessage)
type OnBinMessage = func(belowNexus Nexus, msg []byte)

type Nexus interface {
	// get the id for a given nexus
	UUID() uuid.UUID
	// wether or not two nexuses are equal
	IsEqual(n Nexus) bool
	// called on a message, set in constructor
	OnMessage(belowNexus Nexus, msg base.MoneysocketMessage)
	// called on a binary message, set in constructor
	OnBinMessage(belowNexus Nexus, msg []byte)
	// list all nexuses
	GetDownwardNexusList() []Nexus
	// send a message
	Send(msg base.MoneysocketMessage) error
	// send a binary message
	SendBin(msg []byte) error
	// initiate a close
	InitiateClose()
	// shared seed
	SharedSeed() *beacon.SharedSeed
	// get name of the nexus
	Name() string
	// set a on message callback.
	// This overrides the default (calling on message on below nexus if below nexus is present)
	SetOnMessage(messageFunc OnMessage)
	// set a callback for a binary message
	SetOnBinMessage(messageBinFunc OnBinMessage)
}
