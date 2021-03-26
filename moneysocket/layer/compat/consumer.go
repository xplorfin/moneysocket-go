package compat

type ConsumeNexusInterface interface {
	RequestInvoice(msats int64, overrideRequestUUID, description string)
	RequestPay(bolt11, overrideRequestUUID string)
}
