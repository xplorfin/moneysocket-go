package integrations

import (
	"fmt"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/stack"
)

// WalletConsumer emulates wallet consumer from bs-demo.
type WalletConsumer struct {
	*stack.OutgoingConsumerStack
	ConsumerBeacon beacon.Beacon
}

// generateNewBeacon creates a new beacon based on host, use.
func generateNewBeacon(host string, useTLS bool, port int) beacon.Beacon {
	res := beacon.NewBeacon()
	loc := location.NewWebsocketLocationPort(host, useTLS, port)
	res.AddLocation(loc)
	return res
}

// NewWalletConsumer creates a WalletConsumer from
// a beacon.Beacon and initializes event handlers.
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
