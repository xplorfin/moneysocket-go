package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// NotifyPong is the pong message.
type NotifyPong struct {
	BaseMoneySocketNotification
}

// NewNotifyPong is the pong message.
func NewNotifyPong(requestUUID string) NotifyPong {
	return NotifyPong{NewBaseMoneySocketNotification(base.NotifyPong, requestUUID)}
}

// ToJSON converts the notify pong into a json payload.
func (n NotifyPong) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&m)
}

// DecodeNotifyPong decodes a notify pong message from a json payload.
func DecodeNotifyPong(payload []byte) (NotifyPong, error) {
	notification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyPong{}, err
	}
	return NotifyPong{notification}, nil
}

var _ MoneysocketNotification = &NotifyPong{}
