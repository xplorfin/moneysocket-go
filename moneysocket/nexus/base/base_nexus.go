package base

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"log"
)

// helper function for when youd don't want to pass a handler

type BaseNexus struct {
	// name of the nexus (stored in base for debugging)
	name         string
	uuid         uuid.UUID
	BelowNexus   *nexus.Nexus
	Layer        layer.Layer
	onMessage    nexus.OnMessage
	onBinMessage nexus.OnBinMessage
}

// statically assert nexus type conformity
var _ nexus.Nexus = &BaseNexus{}

func NewBaseNexus(name string) *BaseNexus {
	return &BaseNexus{
		name: name,
		uuid: uuid.NewV4(),
	}
}

func NewBaseNexusBelow(name string, belowNexus nexus.Nexus) *BaseNexus {
	return &BaseNexus{
		name:       name,
		uuid:       uuid.NewV4(),
		BelowNexus: &belowNexus,
	}
}

func NewBaseNexusFull(name string, belowNexus nexus.Nexus, layer layer.Layer) BaseNexus {
	return BaseNexus{
		name:       name,
		uuid:       uuid.NewV4(),
		BelowNexus: &belowNexus,
		Layer:      layer,
	}
}

func (b *BaseNexus) CheckCrossedNexus(belowNexus nexus.Nexus) {
	if b.IsEqual(belowNexus) {
		log.Printf("below nexus: %s (%s) and current nexus %s (%s) appears to be crossed", belowNexus.Name(), belowNexus.Uuid(), b.Name(), b.Uuid())
		log.Print(b.GetDownwardNexusList())
		panic("crossed nexus?")
	}
}

func (b *BaseNexus) Uuid() uuid.UUID {
	return b.uuid
}

func (b *BaseNexus) Name() string {
	return b.name
}

func (b *BaseNexus) IsEqual(n nexus.Nexus) bool {
	return n.Uuid() == b.Uuid()
}

func (b *BaseNexus) OnMessage(belowNexus nexus.Nexus, msg base.MoneysocketMessage) {
	b.CheckCrossedNexus(belowNexus)
	// default to onmessage
	if b.onMessage != nil {
		b.onMessage(belowNexus, msg)
		return
	}
	//if b.BelowNexus != nil {
	//	(*b.BelowNexus).OnMessage(belowNexus, msg)
	//}
}

func (b *BaseNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	b.CheckCrossedNexus(belowNexus)
	// default to onbinmessage
	if b.onBinMessage != nil {
		b.onBinMessage(belowNexus, msg)
		return
	}
	//if b.BelowNexus != nil {
	//	(*b.BelowNexus).OnBinMessage(belowNexus, msg)
	//}
}

func (b BaseNexus) GetDownwardNexusList() (belowList []nexus.Nexus) {
	if b.BelowNexus != nil {
		belowList = (*b.BelowNexus).GetDownwardNexusList()
		belowList = append(belowList, &b)
	}
	return belowList
}

func (b BaseNexus) Send(msg base.MoneysocketMessage) error {
	if b.BelowNexus != nil {
		return (*b.BelowNexus).Send(msg)
	}
	return nil
}

func (b BaseNexus) SendBin(msg []byte) error {
	if b.BelowNexus != nil {
		return (*b.BelowNexus).SendBin(msg)
	}
	return nil
}

func (b BaseNexus) InitiateClose() {
	if b.BelowNexus != nil {
		(*b.BelowNexus).InitiateClose()
	}
}

func (b *BaseNexus) SharedSeed() *beacon.SharedSeed {
	if b.BelowNexus != nil {
		return (*b.BelowNexus).SharedSeed()
	}
	return nil
}

func (b *BaseNexus) SetOnMessage(messageFunc nexus.OnMessage) {
	b.onMessage = messageFunc
}

func (b *BaseNexus) SetOnBinMessage(messageBinFunc nexus.OnBinMessage) {
	b.onBinMessage = messageBinFunc
}
