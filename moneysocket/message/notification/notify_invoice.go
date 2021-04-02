package notification

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// NotifyInvoice gets a Bolt11 invoice from a payload.
type NotifyInvoice struct {
	BaseMoneySocketNotification
	Bolt11 string
}

// NewNotifyInvoice creates a NotifyInvoice bolt11/requestUUID.
func NewNotifyInvoice(bolt11, requestUUID string) NotifyInvoice {
	return NotifyInvoice{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyInvoiceNotification, requestUUID),
		Bolt11:                      bolt11,
	}
}

const bolt11Key = "bolt11"

// MustBeClearText says NotifyInvoice can be encrypted.
func (o NotifyInvoice) MustBeClearText() bool {
	return false
}

// ToJSON gets a json payload from NotifyInvoice.
func (o NotifyInvoice) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(o, m)
	if err != nil {
		return nil, err
	}
	m[bolt11Key] = o.Bolt11
	return json.Marshal(&m)
}

// IsValid determines if the NotifyInvoice is valid.
func (o NotifyInvoice) IsValid() (bool, error) {
	if len(o.Bolt11) < 4 {
		return false, fmt.Errorf("must be bult11")
	}
	if o.Bolt11[0:4] == "lnbc" { // unknown bolt11 type
		return false, fmt.Errorf("doesn't look like a bolt11")
	}
	return true, nil
}

// DecodeNotifyInvoice gets a NotifyInvoice from a payload.
func DecodeNotifyInvoice(payload []byte) (NotifyInvoice, error) {
	notification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyInvoice{}, err
	}
	bolt11Invoice, err := jsonparser.GetString(payload, bolt11Key)
	if err != nil {
		return NotifyInvoice{}, err
	}
	return NotifyInvoice{
		BaseMoneySocketNotification: notification,
		Bolt11:                      bolt11Invoice,
	}, nil
}

var _ MoneysocketNotification = &NotifyInvoice{}
