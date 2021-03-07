package integrations

import (
	"fmt"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/stack"
)

// emulates wallet consumer from bs-demo
type WalletConsumer struct {
	*stack.OutgoingConsumerStack
	ConsumerBeacon beacon.Beacon
}

func makeConsumerBeacon(host string, useTls bool, port int) beacon.Beacon {
	res := beacon.NewBeacon()
	loc := location.NewWebsocketLocationPort(host, false, port)
	res.AddLocation(loc)
	return res
}

func NewWalletConsumer(beacon beacon.Beacon) WalletConsumer {
	cons := stack.NewOutgoingConsumerStack()
	cons.SetOnAnnounce(func(nexus nexus.Nexus) {
		fmt.Println("wallet online")
	})
	cons.SetOnRevoke(func(nexus nexus.Nexus) {
		fmt.Println("wallet offline")
	})
	cons.SetOnProviderInfo(func(consumerTransactNexus nexus.Nexus, msg base.MoneysocketMessage) {
		fmt.Println("provider info")
	})
	cons.SetSendStackEvent(func(layerName string, nexus nexus.Nexus, event string) {
		fmt.Println("stack event")
	})
	cons.SetOnPing(func(nexus nexus.Nexus, msecs int) {
		fmt.Println("pinged")
	})
	cons.SetOnInvoice(func(transactNexus nexus.Nexus, invoice string, requestReferenceUuid string) {
		fmt.Println(invoice)
	})
	cons.SetOnPreimage(func(transactNexus nexus.Nexus, preimage string, requestReferenceUuid string) {
		fmt.Println(preimage)
	})
	return WalletConsumer{cons, beacon}
}
