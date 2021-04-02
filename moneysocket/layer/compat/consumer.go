package compat

// ConsumeNexusInterface is an interface that allows for requesting/paying invoices
type ConsumeNexusInterface interface {
	// RequestInvoice gets an invoice from the lightning driver from msats
	RequestInvoice(msats int64, overrideRequestUUID, description string)
	// RequestPay pays a bolt 11 invoice
	RequestPay(bolt11, overrideRequestUUID string)
}
