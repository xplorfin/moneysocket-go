package local

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	base_moneysocket "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

type Nexus struct {
	*base.NexusBase
}

const LocalNexusName = "LocalNexusName"

func NewLocalNexus(belowNexus nexus.Nexus, layer layer.Layer) *Nexus {
	bnf := base.NewBaseNexusFull(LocalNexusName, belowNexus, layer)
	ln := Nexus{
		&bnf,
	}
	belowNexus.SetOnBinMessage(ln.OnBinMessage)
	belowNexus.SetOnMessage(ln.OnMessage)

	return &ln
}

func (l *Nexus) OnMessage(belowNexus nexus.Nexus, msg base_moneysocket.MoneysocketMessage) {
	log.Println("local nexus got msg")
	l.NexusBase.OnMessage(belowNexus, msg)
}

func (l *Nexus) SendBin(msg []byte) error {
	log.Println("local nexus sent bin")
	return l.NexusBase.SendBin(msg)
}

var _ nexus.Nexus = &Nexus{}
