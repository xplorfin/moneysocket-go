package layer

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

// a non-existent nexus that allows us to avoid returning a null value
type UnknownNexus struct{}

func NewUnknownNexus() UnknownNexus {
	return UnknownNexus{}
}

// TODO disable
func (u UnknownNexus) Uuid() uuid.UUID {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) IsEqual(n nexus.Nexus) bool {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) OnMessage(belowNexus nexus.Nexus, msg base.MoneysocketMessage) {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) GetDownwardNexusList() []nexus.Nexus {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) Send(msg base.MoneysocketMessage) error {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) SendBin(msg []byte) error {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) InitiateClose() {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) SharedSeed() *beacon.SharedSeed {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) Name() string {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) SetOnMessage(messageFunc nexus.OnMessage) {
	panic("this is not a real nexus, did a getter return an error?")
}

func (u UnknownNexus) SetOnBinMessage(messageBinFunc nexus.OnBinMessage) {
	panic("this is not a real nexus, did a getter return an error?")
}


var _ nexus.Nexus = UnknownNexus{}
