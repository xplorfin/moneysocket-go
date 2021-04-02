package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// NotifyOpinionSellerNotReady is.
type NotifyOpinionSellerNotReady struct {
	BaseMoneySocketNotification
}

// NewNotifyOpinionSellerNotReady create a new NotifyOpinionSellerNotReady end notification with a given rendezvous id.
func NewNotifyOpinionSellerNotReady(requestUUID string) NotifyOpinionSellerNotReady {
	return NotifyOpinionSellerNotReady{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyOpinionSellerNotReady, requestUUID),
	}
}

// MustBeClearText determines a NotifyOpinionSellerNotReady.
func (r NotifyOpinionSellerNotReady) MustBeClearText() bool {
	return true
}

// ToJSON marshals a json payload from NotifyOpinionSellerNotReady.
func (r NotifyOpinionSellerNotReady) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&m)
}

// DecodeNotifyOpinionSellerNotReady creates a NotifyOpinionSellerNotReady from a payload.
func DecodeNotifyOpinionSellerNotReady(payload []byte) (NotifyOpinionSellerNotReady, error) {
	notiification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyOpinionSellerNotReady{}, err
	}
	return NotifyOpinionSellerNotReady{notiification}, nil
}

var _ MoneysocketNotification = &NotifyOpinionSellerNotReady{}
