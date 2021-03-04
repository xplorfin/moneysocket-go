package compat

type ConsumeNexusInterface interface {
	RequestInvoice(msats int64, overrideRequestUuid, description string)
	RequestPay(bolt11, overrideRequestUuid string)
}
