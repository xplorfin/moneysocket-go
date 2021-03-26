package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RendezvousEnd struct {
	BaseMoneySocketNotification
	rendezvousID string
}

// create a new rendezvous end notification with a given rendezvous id
func NewRendezvousEnd(rid, requestUUID string) RendezvousEnd {
	return RendezvousEnd{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvousEndNotification, requestUUID),
		rendezvousID:                rid,
	}
}

func (r RendezvousEnd) MustBeClearText() bool {
	return true
}

const rendezvousIDKey = "rendezvous_id"

func (r RendezvousEnd) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIDKey] = r.rendezvousID
	return json.Marshal(&m)
}

func DecodeRendezvousEnd(payload []byte) (RendezvousEnd, error) {
	notiification, err := DecodeRequest(payload)
	if err != nil {
		return RendezvousEnd{}, err
	}
	rendezvousID, err := jsonparser.GetString(payload, rendezvousIDKey)
	if err != nil {
		return RendezvousEnd{}, err
	}
	return RendezvousEnd{notiification, rendezvousID}, nil
}

var _ MoneysocketNotification = &RendezvousEnd{}
