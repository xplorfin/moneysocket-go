package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type NotifyPong struct {
	BaseMoneySocketNotification
}

func NewNotifyPong(requestUUID string) NotifyPong {
	return NotifyPong{NewBaseMoneySocketNotification(base.NotifyPong, requestUUID)}
}

func (n NotifyPong) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&m)
}

func DecodeNotifyPong(payload []byte) (NotifyPong, error) {
	notification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyPong{}, err
	}
	return NotifyPong{notification}, nil
}

var _ MoneysocketNotification = &NotifyPong{}
