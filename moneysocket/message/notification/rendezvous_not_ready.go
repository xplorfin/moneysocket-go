package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RendezvousNotReady struct {
	BaseMoneySocketNotification
	rendezvousID string
}

// create a new rendezvous end notification with a given rendezvous id
func NewRendezvousNotReady(rid, requestUUID string) RendezvousNotReady {
	return RendezvousNotReady{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvousNotReadyNotification, requestUUID),
		rendezvousID:                rid,
	}
}

func (r RendezvousNotReady) MustBeClearText() bool {
	return true
}

func (r RendezvousNotReady) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIDKey] = r.rendezvousID
	return json.Marshal(&m)
}

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
