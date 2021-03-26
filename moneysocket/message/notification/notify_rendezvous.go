package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type Rendezvous struct {
	BaseMoneySocketNotification
	rendezvousID string
}

// create a new rendezvous end notification with a given rendezvous id
func NewNotifyRendezvous(rid, requestUUID string) Rendezvous {
	return Rendezvous{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvous, requestUUID),
		rendezvousID:                rid,
	}
}

func (r Rendezvous) MustBeClearText() bool {
	return true
}

func (r Rendezvous) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIDKey] = r.rendezvousID
	return json.Marshal(&m)
}

func DecodeRendezvous(payload []byte) (Rendezvous, error) {
	notiification, err := DecodeRequest(payload)
	if err != nil {
		return Rendezvous{}, err
	}
	rendezvousID, err := jsonparser.GetString(payload, rendezvousIDKey)
	if err != nil {
		return Rendezvous{}, err
	}
	return Rendezvous{notiification, rendezvousID}, nil
}

var _ MoneysocketNotification = &Rendezvous{}
