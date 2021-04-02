package nexus

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// OnMessage is a handler function for processing a base.MoneysocketMessage
type OnMessage = func(belowNexus Nexus, msg base.MoneysocketMessage)

// OnBinMessage is a handler function for processing a binary message
type OnBinMessage = func(belowNexus Nexus, msg []byte)

// Nexus is the nexus interface
type Nexus interface {
	// UUID gets the id for a given nexus
	UUID() uuid.UUID
	// IsEqual is wether or not two nexuses are equal
	IsEqual(n Nexus) bool
	// OnMessage is called on a message, set in constructor
	OnMessage(belowNexus Nexus, msg base.MoneysocketMessage)
	// OnBinMessage is called on a binary message, set in constructor
	OnBinMessage(belowNexus Nexus, msg []byte)
	// GetDownwardNexusList lists all nexuses
	GetDownwardNexusList() []Nexus
	// Send sends a message
	Send(msg base.MoneysocketMessage) error
	// SendBin sends a binary message
	SendBin(msg []byte) error
	// InitiateClose initiates a close
	InitiateClose()
	// SharedSeed gets a shared seed
	SharedSeed() *beacon.SharedSeed
	// get name of the nexus
	Name() string
	// set a on message callback.
	// This overrides the default (calling on message on below nexus if below nexus is present)
	SetOnMessage(messageFunc OnMessage)
	// set a callback for a binary message
	SetOnBinMessage(messageBinFunc OnBinMessage)
}
