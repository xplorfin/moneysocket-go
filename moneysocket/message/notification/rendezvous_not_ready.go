package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// RendezvousNotReady is a message saying a given rendezvous is not ready.
type RendezvousNotReady struct {
	BaseMoneySocketNotification
	rendezvousID string
}

// NewRendezvousNotReady creates a new rendezvous end notification with a given rendezvous id.
func NewRendezvousNotReady(rid, requestUUID string) RendezvousNotReady {
	return RendezvousNotReady{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvousNotReadyNotification, requestUUID),
		rendezvousID:                rid,
	}
}

// MustBeClearText is wether or not the message is clear.
func (r RendezvousNotReady) MustBeClearText() bool {
	return true
}

// ToJSON converts the message to json.
func (r RendezvousNotReady) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIDKey] = r.rendezvousID
	return json.Marshal(&m)
}

// DecodeRendezvousNotReady converts the payload to RendezvousNotReady.
func DecodeRendezvousNotReady(payload []byte) (RendezvousNotReady, error) {
	notiification, err := DecodeRequest(payload)
	if err != nil {
		return RendezvousNotReady{}, err
	}
	rendezvousID, err := jsonparser.GetString(payload, rendezvousIDKey)
	if err != nil {
		return RendezvousNotReady{}, err
	}
	return RendezvousNotReady{notiification, rendezvousID}, nil
}

var _ MoneysocketNotification = &RendezvousNotReady{}
