package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// Rendezvous notifies that a rendezvous is ready by id.
type Rendezvous struct {
	BaseMoneySocketNotification
	rendezvousID string
}

// NewNotifyRendezvous creates a new rendezvous end notification with a given rendezvous id.
func NewNotifyRendezvous(rid, requestUUID string) Rendezvous {
	return Rendezvous{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvous, requestUUID),
		rendezvousID:                rid,
	}
}

// MustBeClearText text denotes a Rendezvous can be clear text.
func (r Rendezvous) MustBeClearText() bool {
	return true
}

// ToJSON converts a Rendezvous notification to json.
func (r Rendezvous) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIDKey] = r.rendezvousID
	return json.Marshal(&m)
}

// DecodeRendezvous decodes a Rendezvous notification from json.
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
