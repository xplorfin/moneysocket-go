package notification

import (
	"fmt"

	"github.com/buger/jsonparser"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

func NotificationFromText(payload []byte) (base.MoneysocketMessage, base.MessageType, error) {
	rawType, err := jsonparser.GetString(payload, NotificationNameKey)
	if err != nil {
		return nil, 0, err
	}
	msgType := base.MessageTypeFromString(rawType)
	switch msgType {
	case base.NotifyRendezvous:
		decoded, err := DecodeRendezvous(payload)
		return decoded, base.NotifyRendezvous, err
	case base.NotifyRendezvousNotReadyNotification:
		decoded, err := DecodeRendezvousNotReady(payload)
		return decoded, base.NotifyRendezvousNotReadyNotification, err
	case base.NotifyRendezvousEndNotification:
		decoded, err := DecodeRendezvousEnd(payload)
		return decoded, base.NotifyRendezvousEndNotification, err
	case base.NotifyInvoiceNotification:
		decoded, err := DecodeNotifyInvoice(payload)
		return decoded, base.NotifyInvoiceNotification, err
	case base.NotifyOpinionSeller:
		decoded, err := DecodeNotifyOpinionSeller(payload)
		return decoded, base.NotifyOpinionSeller, err
	case base.NotifyOpinionSellerNotReady:
		decoded, err := DecodeNotifyOpinionSellerNotReady(payload)
		return decoded, base.NotifyOpinionSellerNotReady, err
	case base.NotifyOpinionInvoice:
		decoded, err := DecodeNotifyOpinionInvoice(payload)
		return decoded, base.NotifyOpinionInvoice, err
	case base.NotifyProviderNotReady:
		decoded, err := DecodeNotifyProviderNotReady(payload)
		return decoded, base.NotifyProviderNotReady, err
	case base.NotifyPong:
		decoded, err := DecodeNotifyPong(payload)
		return decoded, base.NotifyPong, err
	}
	return nil, 0, fmt.Errorf("message type %s not yet implemented", msgType.ToString())
}
