package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type NotifyOpinionInvoice struct {
	BaseMoneySocketNotification
	Bolt11 string
}

func NewNotifyOpinionInvoice(requestUuid, bolt11 string) NotifyOpinionInvoice {
	return NotifyOpinionInvoice{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyOpinionInvoice, requestUuid),
		Bolt11:                      bolt11,
	}
}

func (n NotifyOpinionInvoice) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	m[bolt11Key] = n.Bolt11
	return json.Marshal(&m)
}

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
