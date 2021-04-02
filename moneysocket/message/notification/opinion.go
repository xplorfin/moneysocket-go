package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// NotifyOpinionInvoice message.
type NotifyOpinionInvoice struct {
	BaseMoneySocketNotification
	// Bolt11 is the bolt11 invoice
	Bolt11 string
}

// NewNotifyOpinionInvoice creates a new NotifyOpinionInvoice message.
func NewNotifyOpinionInvoice(requestUUID, bolt11 string) NotifyOpinionInvoice {
	return NotifyOpinionInvoice{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyOpinionInvoice, requestUUID),
		Bolt11:                      bolt11,
	}
}

// ToJSON converts a NotifyOpinionInvoice to a json payload.
func (n NotifyOpinionInvoice) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	m[bolt11Key] = n.Bolt11
	return json.Marshal(&m)
}

// DecodeNotifyOpinionInvoice creates a new NotifyOpinionInvoice from a payload.
func DecodeNotifyOpinionInvoice(payload []byte) (NotifyOpinionInvoice, error) {
	notification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyOpinionInvoice{}, err
	}
	bolt11Invoice, err := jsonparser.GetString(payload, bolt11Key)
	if err != nil {
		return NotifyOpinionInvoice{}, err
	}
	return NotifyOpinionInvoice{
		BaseMoneySocketNotification: notification,
		Bolt11:                      bolt11Invoice,
	}, nil
}

var _ MoneysocketNotification = &NotifyOpinionInvoice{}
