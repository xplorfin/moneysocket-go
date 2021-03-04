package integrations

import (
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/seller"
	"github.com/xplorfin/moneysocket-go/moneysocket/stack"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

type SellerApp struct {
	SellerStack    *stack.SellerStack
	ConsumerStack  *stack.OutgoingConsumerStack
	AccountUuid    string
	SellerUuid     string
	RequestedItems map[string]string
	StoreOpen      bool
	ProviderInfo   seller.SellerInfo
}

func NewSellerApp(host string, useTls bool, port int) *SellerApp {
	walletConsumer := NewWalletConsumer(host, useTls, port)
	sellerStack := stack.NewSellerStack()
	sa := SellerApp{
		SellerStack:    sellerStack,
		ConsumerStack:  walletConsumer.OutgoingConsumerStack,
		AccountUuid:    uuid.NewV4().String(),
		SellerUuid:     uuid.NewV4().String(),
		RequestedItems: make(map[string]string),
	}
	return &sa
}

func (sa *SellerApp) OpenStore() {
	sa.StoreOpen = true
	sa.SellerStack.SellerNowReadyFromApp()
}

func (sa *SellerApp) CloseStore() {
	sa.StoreOpen = false
	sa.SellerStack.DoDisconnect()
}

func (sa *SellerApp) UpdatePrices() {
	sa.SellerStack.UpdatePrices()
}

const (
	// names
	helloName   = "hello"
	timeName    = "time"
	outlookName = "outlook"
	// prices
	helloPrice   = 50
	timePrice    = 100
	outlookPrice = 150
	// description
	helloDescription         = "Hello World"
	timestampDescription     = "Current Timestamp"
	marketOutlookDescription = "Market Outlook"
)

func (sa *SellerApp) SetupSellerStack() {
	sa.SellerStack.SetOnAnnounce(func(nexus nexus.Nexus) {
		log.Println("provider online")
	})
	sa.SellerStack.SetOnRevoke(func(nexus nexus.Nexus) {
		log.Println("nexus revoked (offline)")
	})
	sa.SellerStack.SetOnStackEvent(func(layerName string, nexus nexus.Nexus, status string) {
		log.Println(status)
	})
	sa.SellerStack.SetHandleOpinionInvoiceRequest(func(item string, requestUuid string) {
		switch item {
		case helloName:
			sa.ConsumerStack.RequestInvoice(helloPrice, requestUuid, helloDescription)
		case timeName:
			sa.ConsumerStack.RequestInvoice(timePrice, requestUuid, timestampDescription)
		case outlookName:
			sa.ConsumerStack.RequestInvoice(outlookPrice, requestUuid, marketOutlookDescription)
		default:
			panic(fmt.Errorf("unknown item id"))
		}
		log.Println(fmt.Sprintf("adding invoice request: %s", requestUuid))
	})
	sa.SellerStack.SetHandleSellerInfoRequest(func() seller.SellerInfo {
		if sa.StoreOpen {
			return seller.SellerInfo{
				Ready:      true,
				SellerUUID: sa.SellerUuid,
				Items: []notification.Item{
					{
						ItemID: helloName,
						Name:   helloDescription,
						Msats:  helloPrice,
					},
					{
						ItemID: timeName,
						Name:   timestampDescription,
						Msats:  timePrice,
					},
					{
						ItemID: outlookName,
						Name:   marketOutlookDescription,
						Msats:  outlookPrice,
					},
				},
			}
		} else {
			return seller.SellerInfo{Ready: false}
		}
	})
	sa.SellerStack.SetHandleProviderInfoRequest(func(seed beacon.SharedSeed) account.AccountDb {
		panic("TODO")
	})
}
