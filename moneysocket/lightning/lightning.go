package lightning

// TODO this is a callback
type PaidCallback func(preimage string, msats int)

// Lightning is an implementable interface for interacting with different lightning clients
type Lightning interface {
	// RegisterPaidRecvCb registers a PaidCallback for when an invoice is paid
	RegisterPaidRecvCb(callback PaidCallback)
	// GetInvoice retrieves a bolt11 invoice for a given msatAmount
	GetInvoice(msatAmount int) (paymentRequest string, err error)
	// PayInvoice pays a bolt11 invoice
	PayInvoice(bolt11 string) (preimage []byte, msatAmount int, err error)
	// RecvPaid receives a payment
	RecvPaid(preimage string, msats int)
}
