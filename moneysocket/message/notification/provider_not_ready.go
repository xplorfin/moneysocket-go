package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// NotifyProviderNotReady notifies a provider is not messaged.
type NotifyProviderNotReady struct {
	BaseMoneySocketNotification
}

// NewNotifyProviderNotReady creates a message that a provider is not ready.
func NewNotifyProviderNotReady(requestUUID string) NotifyProviderNotReady {
	return NotifyProviderNotReady{NewBaseMoneySocketNotification(base.NotifyProviderNotReady, requestUUID)}
}

// ToJSON marshalls a NotifyProviderNotReady message into a json payload.
func (n NotifyProviderNotReady) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&m)
}

// DecodeNotifyProviderNotReady unmarshalls json into a NotifyProviderNotReady message.
func DecodeNotifyProviderNotReady(payload []byte) (NotifyProviderNotReady, error) {
	notification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyProviderNotReady{}, err
	}
	return NotifyProviderNotReady{notification}, nil
}

var _ MoneysocketNotification = &NotifyProviderNotReady{}
