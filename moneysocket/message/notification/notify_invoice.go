package notification

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type NotifyInvoice struct {
	BaseMoneySocketNotification
	Bolt11 string
}

func NewNotifyInvoice(bolt11, requestUuid string) NotifyInvoice {
	return NotifyInvoice{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyInvoiceNotification, requestUuid),
		Bolt11:                      bolt11,
	}
}

const bolt11Key = "bolt11"

func (o NotifyInvoice) MustBeClearText() bool {
	return false
}

func (o NotifyInvoice) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(o, m)
	if err != nil {
		return nil, err
	}
	m[bolt11Key] = o.Bolt11
	return json.Marshal(&m)
}

func (o NotifyInvoice) IsValid() (bool, error) {
	if len(o.Bolt11) < 4 {
		return false, fmt.Errorf("must be bult11")
	}
	if o.Bolt11[0:4] == "lnbc" { // unknown bolt11 type
		return false, fmt.Errorf("doesn't look like a bolt11")
	}
	return true, nil
}

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
