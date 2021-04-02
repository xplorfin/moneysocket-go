package base

import (
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// helper function for when youd don't want to pass a handler

// NexusBase is the nexus superclass. It contains common functions for a nexus
type NexusBase struct {
	// name of the nexus (stored in base for debugging)
	name         string
	uuid         uuid.UUID
	BelowNexus   *nexus.Nexus
	Layer        layer.Base
	onMessage    nexus.OnMessage
	onBinMessage nexus.OnBinMessage
}

// statically assert nexus type conformity
var _ nexus.Nexus = &NexusBase{}

// NewBaseNexus creates a new nexus base
func NewBaseNexus(name string) *NexusBase {
	return &NexusBase{
		name: name,
		uuid: uuid.NewV4(),
	}
}

// NewBaseNexusBelow creates a new base nexus and sets the below nexus
func NewBaseNexusBelow(name string, belowNexus nexus.Nexus) *NexusBase {
	return &NexusBase{
		name:       name,
		uuid:       uuid.NewV4(),
		BelowNexus: &belowNexus,
	}
}

// NewBaseNexusFull creates a new base nexus and sets a below nexus/layer for comms
func NewBaseNexusFull(name string, belowNexus nexus.Nexus, layer layer.Base) NexusBase {
	return NexusBase{
		name:       name,
		uuid:       uuid.NewV4(),
		BelowNexus: &belowNexus,
		Layer:      layer,
	}
}

// CheckCrossedNexus checks if the nexus has been crossed
func (b *NexusBase) CheckCrossedNexus(belowNexus nexus.Nexus) {
	if b.IsEqual(belowNexus) {
		log.Printf("below nexus: %s (%s) and current nexus %s (%s) appears to be crossed", belowNexus.Name(), belowNexus.UUID(), b.Name(), b.UUID())
		log.Print(b.GetDownwardNexusList())
		panic("crossed nexus?")
	}
}

// UUID gets the uuid of the nexus base
func (b *NexusBase) UUID() uuid.UUID {
	return b.uuid
}

// Name gets the name of the nexus
func (b *NexusBase) Name() string {
	return b.name
}

// IsEqual checks if a nexus is equal to another nexus
func (b *NexusBase) IsEqual(n nexus.Nexus) bool {
	return n.UUID() == b.UUID()
}

// OnMessage handles the message callback
func (b *NexusBase) OnMessage(belowNexus nexus.Nexus, msg base.MoneysocketMessage) {
	b.CheckCrossedNexus(belowNexus)
	if b.onMessage != nil {
		b.onMessage(b, msg)
		return
	}
}

// OnBinMessage calls the binary message method
func (b *NexusBase) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	b.CheckCrossedNexus(belowNexus)
	// default to onbinmessage
	if b.onBinMessage != nil {
		b.onBinMessage(b, msg)
		return
	}
}

// GetDownwardNexusList gets the nexus list
func (b NexusBase) GetDownwardNexusList() (belowList []nexus.Nexus) {
	if b.BelowNexus != nil {
		belowList = (*b.BelowNexus).GetDownwardNexusList()
		belowList = append(belowList, &b)
	}
	return belowList
}

// Send sends a message
func (b *NexusBase) Send(msg base.MoneysocketMessage) error {
	if b.BelowNexus != nil {
		return (*b.BelowNexus).Send(msg)
	}
	return nil
}

// SendBin sends a binary message
func (b *NexusBase) SendBin(msg []byte) error {
	if b.BelowNexus != nil {
		return (*b.BelowNexus).SendBin(msg)
	}
	return nil
}

// InitiateClose closes the websocket
func (b *NexusBase) InitiateClose() {
	if b.BelowNexus != nil {
		(*b.BelowNexus).InitiateClose()
	}
}

// SharedSeed gets the shared seed
func (b NexusBase) SharedSeed() *beacon.SharedSeed {
	if b.BelowNexus != nil {
		return (*b.BelowNexus).SharedSeed()
	}
	return nil
}

// SetOnMessage sets the message callback function
func (b *NexusBase) SetOnMessage(messageFunc nexus.OnMessage) {
	b.onMessage = messageFunc
}

// SetOnBinMessage sets the binary message callback function
func (b *NexusBase) SetOnBinMessage(messageBinFunc nexus.OnBinMessage) {
	b.onBinMessage = messageBinFunc
}
