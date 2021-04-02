package compat

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

// HandleInvoiceRequest processes an invoice request for msats
type HandleInvoiceRequest func(nexus nexus.Nexus, msats int64, requestUUID string)

// HandlePayRequest pays a bolt11 invoice
type HandlePayRequest func(nexus nexus.Nexus, bolt11 string, requestUUID string)

// HandleProviderInfoRequest is a sesed
type HandleProviderInfoRequest func(seed beacon.SharedSeed) account.DB

// HandleOpinionInvoiceRequest pays a request for an item
type HandleOpinionInvoiceRequest func(nx nexus.Nexus, itemId string, requestUUID string)

// ProviderTransactLayerInterface is a provider for handling provider requests
type ProviderTransactLayerInterface interface {
	layer.Base
	// HandleProviderInfoRequest processes an info request with a seed
	HandleProviderInfoRequest(seed beacon.SharedSeed) account.DB
	// HandlePayRequest processes a payment request
	HandlePayRequest(nexus nexus.Nexus, bolt11 string, requestUUID string)
	// HandleInvoiceRequest processes an invoice request
	HandleInvoiceRequest(nexus nexus.Nexus, msats int64, requestUUID string)
}

// SellingLayerInterface is an interface with a NexusWaitingForApp method
type SellingLayerInterface interface {
	// NexusWaitingForApp registers the nexus with the NexusWaitingForApp method
	NexusWaitingForApp(seed *beacon.SharedSeed, sellerNexus nexus.Nexus)
}

// WaitingForApp is a map of uuid->nexus.Nexus
type WaitingForApp map[string]nexus.Nexus
