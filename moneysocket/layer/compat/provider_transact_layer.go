package compat

import (
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

type HandleInvoiceRequest func(nexus nexus.Nexus, msats int64, requestUuid string)
type HandlePayRequest func(nexus nexus.Nexus, bolt11 string, requestUuid string)
type HandleProviderInfoRequest func(seed beacon.SharedSeed) account.AccountDb
type HandleOpinionInvoiceRequest func(nx nexus.Nexus, itemId string, requestUuid string)

type ProviderTransactLayerInterface interface {
	layer.Layer
	HandleProviderInfoRequest(seed beacon.SharedSeed) account.AccountDb
	HandlePayRequest(nexus nexus.Nexus, bolt11 string, requestUuid string)
	HandleInvoiceRequest(nexus nexus.Nexus, msats int64, requestUuid string)
}

type SellingLayerInterface interface {
	NexusWaitingForApp(seed *beacon.SharedSeed, sellerNexus nexus.Nexus)
}

type WaitingForApp map[string]nexus.Nexus
