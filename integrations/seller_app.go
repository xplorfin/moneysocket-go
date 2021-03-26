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

// SellerApp is a prototype based on https://git.io/JmVLj designed to test interactions
// between two accounts on terminus
type SellerApp struct {
	// SellerStack is a stack of layers to manage events in the seller provider app https://git.io/JmVq7
	SellerStack *stack.SellerStack
	// ConsumerStack is a stack of layers to manage events in the seller consumer app
	ConsumerStack *stack.OutgoingConsumerStack
	// AccountUUID is the uuid used to identify the (buyer?) account
	AccountUUID string
	// SellerUUID is the uuid used to identify the (seller) account
	SellerUUID string
	// RequestedItems is a list of items the buyer "checks out with"
	RequestedItems map[string]string
	// StoreOpen is whether or not the seller opens the store
	StoreOpen bool
	// ProviderInfo contains information about the state of the sale (who's paying who, the wad, account uuid)
	ProviderInfo seller.Info
}

// NewSellerApp initializes a seller app with a beacon (that contains the location we communicate with)
func NewSellerApp(beacon beacon.Beacon) *SellerApp {
	walletConsumer := NewWalletConsumer(beacon)
	sellerStack := stack.NewSellerStack()
	sa := SellerApp{
		SellerStack:    sellerStack,
		ConsumerStack:  walletConsumer.OutgoingConsumerStack,
		AccountUUID:    uuid.NewV4().String(),
		SellerUUID:     uuid.NewV4().String(),
		RequestedItems: make(map[string]string),
	}
	return &sa
}

// OpenStore opens the shop
func (sa *SellerApp) OpenStore() {
	sa.StoreOpen = true
	sa.SellerStack.SellerNowReadyFromApp()
}

// CloseStore closes the shop
func (sa *SellerApp) CloseStore() {
	sa.StoreOpen = false
	sa.SellerStack.DoDisconnect()
}

// UpdatePrices updates the prices in the state
func (sa *SellerApp) UpdatePrices() {
	sa.SellerStack.UpdatePrices()
}

// items names/prices/descriptions
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

// SetupSellerStack creates eevent handlers on the seller stack for interacting w/ the consuer stack
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
	sa.SellerStack.SetHandleSellerInfoRequest(func() seller.Info {
		if sa.StoreOpen {
			return seller.Info{
				Ready:      true,
				SellerUUID: sa.SellerUUID,
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
		}
		return seller.Info{Ready: false}
	})
	sa.SellerStack.SetHandleProviderInfoRequest(func(seed beacon.SharedSeed) account.Db {
		panic("TODO")
	})
}
