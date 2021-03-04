package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type NotifyProviderNotReady struct {
	BaseMoneySocketNotification
}

func NewNotifyProviderNotReady(requestUuid string) NotifyProviderNotReady {
	return NotifyProviderNotReady{NewBaseMoneySocketNotification(base.NotifyProviderNotReady, requestUuid)}
}

func (n NotifyProviderNotReady) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&m)
}

func DecodeNotifyProviderNotReady(payload []byte) (NotifyProviderNotReady, error) {
	notification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyProviderNotReady{}, err
	}
	return NotifyProviderNotReady{notification}, nil
}

var _ MoneysocketNotification = &NotifyProviderNotReady{}
