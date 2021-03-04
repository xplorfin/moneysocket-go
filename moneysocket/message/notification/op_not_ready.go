package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type NotifyOpinionSellerNotReady struct {
	BaseMoneySocketNotification
}

// create a new rendezvous end notification with a given rendezvous id
func NewNotifyOpinionSellerNotReady(requestUuid string) NotifyOpinionSellerNotReady {
	return NotifyOpinionSellerNotReady{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyOpinionSellerNotReady, requestUuid),
	}
}

func (r NotifyOpinionSellerNotReady) MustBeClearText() bool {
	return true
}

func (r NotifyOpinionSellerNotReady) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&m)
}

func DecodeNotifyOpinionSellerNotReady(payload []byte) (NotifyOpinionSellerNotReady, error) {
	notiification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyOpinionSellerNotReady{}, err
	}
	return NotifyOpinionSellerNotReady{notiification}, nil
}

var _ MoneysocketNotification = &NotifyOpinionSellerNotReady{}
