package base

//covers both request and notifications
type MessageType int

const (
	// requests
	PingRequest           MessageType = 0
	ProviderRequest       MessageType = iota
	InvoiceRequest        MessageType = iota
	PayRequest            MessageType = iota
	RendezvousRequest     MessageType = iota
	RequestOpinionSeller  MessageType = iota
	RequestOpinionInvoice MessageType = iota

	// notifications
	NotifyRendezvous                     = iota
	NotifyRendezvousEndNotification      = iota
	NotifyRendezvousNotReadyNotification = iota
	NotifyInvoiceNotification            = iota
	NotifyPreimage                       = iota
	NotifyProvider                       = iota
	NotifyOpinionSeller                  = iota
	NotifyOpinionSellerNotReady          = iota
	NotifyProviderNotReady               = iota
	NotifyOpinionInvoice                 = iota
	NotifyPong                           = iota
)

// used for
const (
	// requests
	RequestPingName           = "REQUEST_PING"
	RequestProviderName       = "REQUEST_PROVIDER"
	RequestInvoiceName        = "REQUEST_INVOICE"
	RequestPayName            = "REQUEST_PAY"
	RendezvousRequestName     = "REQUEST_RENDEZVOUS"
	RequestOpinionSellerName  = "REQUEST_OPINION_SELLER"
	RequestOpinionInvoiceName = "REQUEST_OPINION_INVOICE"

	// notifications
	NotifyNotifyRendezvousEndName   = "NOTIFY_RENDEZVOUS_END"
	NotifyInvoiceName               = "NOTIFY_INVOICE"
	NotifyPreimageName              = "NOTIFY_PREIMAGE"
	NotifyProviderName              = "NOTIFY_PROVIDER"
	NotifyOpinionSellerName         = "NOTIFY_OPINION_SELLER"
	NotifyOpinionSellerNotReadyName = "NOTIFY_OPINION_SELLER_NOT_READY"
	NotifyOpinionInvoiceName        = "NOTIFY_OPINION_INVOICE"
	NotifyProviderNotReadyName      = "NOTIFY_PROVIDER_NOT_READY"
	NotifyPongName                  = "NOTIFY_PONG"
	NotifyRendezvousName            = "NOTIFY_RENDEZVOUS"
	NotifyRendezvousNotReadyName    = "NOTIFY_RENDEZVOUS_NOT_READY"
)

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
