package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// RendezvousEnd ends a rendezvous of a given id.
type RendezvousEnd struct {
	BaseMoneySocketNotification
	rendezvousID string
}

// NewRendezvousEnd creates a new rendezvous end notification with a given rendezvous id.
func NewRendezvousEnd(rid, requestUUID string) RendezvousEnd {
	return RendezvousEnd{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvousEndNotification, requestUUID),
		rendezvousID:                rid,
	}
}

// MustBeClearText determines wether or not the message must be clear text.
func (r RendezvousEnd) MustBeClearText() bool {
	return true
}

const rendezvousIDKey = "rendezvous_id"

// ToJSON converts a RendezvousEnd message to a json payload.
func (r RendezvousEnd) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIDKey] = r.rendezvousID
	return json.Marshal(&m)
}

// DecodeRendezvousEnd decodes a rendezvous message.
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
