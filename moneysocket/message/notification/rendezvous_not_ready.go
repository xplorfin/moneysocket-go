package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RendezvousNotReady struct {
	BaseMoneySocketNotification
	rendezvousId string
}

// create a new rendezvous end notification with a given rendezvous id
func NewRendezvousNotReady(rid, requestUuid string) RendezvousNotReady {
	return RendezvousNotReady{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvousNotReadyNotification, requestUuid),
		rendezvousId:                rid,
	}
}

func (r RendezvousNotReady) MustBeClearText() bool {
	return true
}

func (r RendezvousNotReady) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIdKey] = r.rendezvousId
	return json.Marshal(&m)
}

func DecodeRendezvousNotReady(payload []byte) (RendezvousNotReady, error) {
	notiification, err := DecodeRequest(payload)
	if err != nil {
		return RendezvousNotReady{}, err
	}
	rendezvousId, err := jsonparser.GetString(payload, rendezvousIdKey)
	if err != nil {
		return RendezvousNotReady{}, err
	}
	return RendezvousNotReady{notiification, rendezvousId}, nil
}

var _ MoneysocketNotification = &RendezvousNotReady{}
