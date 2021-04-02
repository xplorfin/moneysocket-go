package base

// MessageType is the message or notification type.
type MessageType int

const (
	// PingRequest is a ping request.
	PingRequest MessageType = 0
	// ProviderRequest is a request for a provider.
	ProviderRequest MessageType = iota
	// InvoiceRequest is a request for an invoice.
	InvoiceRequest MessageType = iota
	// PayRequest is a request for payment.
	PayRequest MessageType = iota
	// RendezvousRequest is a request for a provider.
	RendezvousRequest MessageType = iota
	// RequestOpinionSeller requests items from an opinion seller.
	RequestOpinionSeller MessageType = iota
	// RequestOpinionInvoice creates an invoice for an opinion.
	RequestOpinionInvoice MessageType = iota

	// NotifyRendezvous notifies a rendezvous is ready.
	NotifyRendezvous = iota
	// NotifyRendezvousEndNotification notifies a rendezvous has ended.
	NotifyRendezvousEndNotification = iota
	// NotifyRendezvousNotReadyNotification notifies a rendezvous is not ready.
	NotifyRendezvousNotReadyNotification = iota
	// NotifyInvoiceNotification notifies a new invoice has come in.
	NotifyInvoiceNotification = iota
	// NotifyPreimage notifies a preimage has come in.
	NotifyPreimage = iota
	// NotifyProvider notifies a provider is ready.
	NotifyProvider = iota
	// NotifyOpinionSeller notifies an opinion seller is ready.
	NotifyOpinionSeller = iota
	// NotifyOpinionSellerNotReady notifies an opinion seller is ready.
	NotifyOpinionSellerNotReady = iota
	// NotifyProviderNotReady notifies a provider is not ready.
	NotifyProviderNotReady = iota
	// NotifyOpinionInvoice notify an opinion invoice.
	NotifyOpinionInvoice = iota
	// NotifyPing notifies a ping.
	NotifyPing = iota
	// NotifyPong notifies a pong response.
	NotifyPong = iota
)

// used for.
const (
	// RequestPingName requests a ping message.
	RequestPingName = "REQUEST_PING"
	// RequestProviderName requests a provider.
	RequestProviderName = "REQUEST_PROVIDER"
	// RequestInvoiceName requests an invoice.
	RequestInvoiceName = "REQUEST_INVOICE"
	// RequestPayName requests a payment.
	RequestPayName = "REQUEST_PAY"
	// RendezvousRequestName requests a rendezvous.
	RendezvousRequestName = "REQUEST_RENDEZVOUS"
	// RequestOpinionSellerName requests an opinion seller.
	RequestOpinionSellerName = "REQUEST_OPINION_SELLER"
	// RequestOpinionInvoiceName requests an opinion invoice.
	RequestOpinionInvoiceName = "REQUEST_OPINION_INVOICE"

	// NotifyNotifyRendezvousEndName notifies a rendezvous end.
	NotifyNotifyRendezvousEndName = "NOTIFY_RENDEZVOUS_END"
	// NotifyInvoiceName notifies an invoice is ready.
	NotifyInvoiceName = "NOTIFY_INVOICE"
	// NotifyPreimageName notifies a preimage.
	NotifyPreimageName = "NOTIFY_PREIMAGE"
	// NotifyProviderName notifies a provider.
	NotifyProviderName = "NOTIFY_PROVIDER"
	// NotifyOpinionSellerName notifies a seller is ready.
	NotifyOpinionSellerName = "NOTIFY_OPINION_SELLER"
	// NotifyOpinionSellerNotReadyName notifies a seller is not ready.
	NotifyOpinionSellerNotReadyName = "NOTIFY_OPINION_SELLER_NOT_READY"
	// NotifyOpinionInvoiceName notifies an invoice.
	NotifyOpinionInvoiceName = "NOTIFY_OPINION_INVOICE"
	// NotifyProviderNotReadyName notifies a provider is not ready.
	NotifyProviderNotReadyName = "NOTIFY_PROVIDER_NOT_READY"
	// NotifyPongName notifies a pong reply.
	NotifyPongName = "NOTIFY_PONG"
	// NotifyRendezvousName notifies a rendezvous is ready.
	NotifyRendezvousName = "NOTIFY_RENDEZVOUS"
	// NotifyRendezvousNotReadyName notifies a rendezvous is not readt.
	NotifyRendezvousNotReadyName = "NOTIFY_RENDEZVOUS_NOT_READY"
)

// ToString converts a message to a string.
func (r MessageType) ToString() string {
	switch r {
	// requests
	case PingRequest:
		return RequestPingName
	case ProviderRequest:
		return RequestProviderName
	case InvoiceRequest:
		return RequestInvoiceName
	case PayRequest:
		return RequestPayName
	case RendezvousRequest:
		return RendezvousRequestName
	case RequestOpinionInvoice:
		return RequestOpinionInvoiceName
	case RequestOpinionSeller:
		return RequestOpinionSellerName
	// notifications
	case NotifyRendezvousEndNotification:
		return NotifyNotifyRendezvousEndName
	case NotifyInvoiceNotification:
		return NotifyInvoiceName
	case NotifyPreimage:
		return NotifyPreimageName
	case NotifyProvider:
		return NotifyProviderName
	case NotifyOpinionSeller:
		return NotifyOpinionSellerName
	case NotifyOpinionSellerNotReady:
		return NotifyOpinionSellerNotReadyName
	case NotifyOpinionInvoice:
		return NotifyOpinionInvoiceName
	case NotifyProviderNotReady:
		return NotifyProviderNotReadyName
	case NotifyPong:
		return NotifyPongName
	case NotifyRendezvous:
		return NotifyRendezvousName
	case NotifyRendezvousNotReadyNotification:
		return NotifyRendezvousNotReadyName
	}
	panic("message not yet implemented")
}

// MessageTypeFromString converts a name to a MessageType.
func MessageTypeFromString(name string) MessageType {
	switch name {
	// requests
	case RequestPingName:
		return PingRequest
	case RequestProviderName:
		return ProviderRequest
	case RequestInvoiceName:
		return InvoiceRequest
	case RequestPayName:
		return PayRequest
	case RendezvousRequestName:
		return RendezvousRequest
	case RequestOpinionInvoiceName:
		return RequestOpinionInvoice
	case RequestOpinionSellerName:
		return RequestOpinionSeller
	// notifications
	case NotifyNotifyRendezvousEndName:
		return NotifyRendezvousEndNotification
	case NotifyInvoiceName:
		return NotifyInvoiceNotification
	case NotifyPreimageName:
		return NotifyPreimage
	case NotifyProviderName:
		return NotifyProvider
	case NotifyOpinionSellerName:
		return NotifyOpinionSeller
	case NotifyOpinionSellerNotReadyName:
		return NotifyOpinionSellerNotReady
	case NotifyOpinionInvoiceName:
		return NotifyOpinionInvoice
	case NotifyProviderNotReadyName:
		return NotifyProviderNotReady
	case NotifyPongName:
		return NotifyPong
	case NotifyRendezvousName:
		return NotifyRendezvous
	case NotifyRendezvousNotReadyName:
		return NotifyRendezvousNotReadyNotification
	}
	panic("name not found")
}
