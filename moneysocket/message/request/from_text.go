package request

import (
	"fmt"

	"github.com/buger/jsonparser"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// FromText converts a payload to a decoded moneysocket message
// TODO handle more elegantly
func FromText(payload []byte) (base.MoneysocketMessage, base.MessageType, error) {
	rawType, err := jsonparser.GetString(payload, NameKey)
	if err != nil {
		return nil, 0, err
	}
	msgType := base.MessageTypeFromString(rawType)
	switch msgType {
	case base.PingRequest:
		decoded, err := DecodePing(payload)
		return decoded, base.PingRequest, err
	case base.ProviderRequest:
		decoded, err := DecodeRequestProvider(payload)
		return decoded, base.ProviderRequest, err
	case base.InvoiceRequest:
		decoded, err := DecodeRequestInvoice(payload)
		return decoded, base.InvoiceRequest, err
	case base.PayRequest:
		decoded, err := DecodeRequestPay(payload)
		return decoded, base.PayRequest, err
	case base.RendezvousRequest:
		decoded, err := DecodeRendezvousRequest(payload)
		return decoded, base.RendezvousRequest, err
	case base.RequestOpinionInvoice:
		decoded, err := DecodeRequestOpinionInvoice(payload)
		return decoded, base.RequestOpinionInvoice, err
	}

	return nil, 0, fmt.Errorf("message type %s not yet implemented", msgType.ToString())
}
