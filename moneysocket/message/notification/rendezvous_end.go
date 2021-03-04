package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RendezvousEnd struct {
	BaseMoneySocketNotification
	rendezvousId string
}

// create a new rendezvous end notification with a given rendezvous id
func NewRendezvousEnd(rid, requestUuid string) RendezvousEnd {
	return RendezvousEnd{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvousEndNotification, requestUuid),
		rendezvousId:                rid,
	}
}

func (r RendezvousEnd) MustBeClearText() bool {
	return true
}

const rendezvousIdKey = "rendezvous_id"

func (r RendezvousEnd) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIdKey] = r.rendezvousId
	return json.Marshal(&m)
}

func DecodeRendezvousEnd(payload []byte) (RendezvousEnd, error) {
	notiification, err := DecodeRequest(payload)
	if err != nil {
		return RendezvousEnd{}, err
	}
	rendezvousId, err := jsonparser.GetString(payload, rendezvousIdKey)
	if err != nil {
		return RendezvousEnd{}, err
	}
	return RendezvousEnd{notiification, rendezvousId}, nil
}

var _ MoneysocketNotification = &RendezvousEnd{}
