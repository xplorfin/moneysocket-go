package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type Rendezvous struct {
	BaseMoneySocketNotification
	rendezvousId string
}

// create a new rendezvous end notification with a given rendezvous id
func NewNotifyRendezvous(rid, requestUuid string) Rendezvous {
	return Rendezvous{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyRendezvous, requestUuid),
		rendezvousId:                rid,
	}
}

func (r Rendezvous) MustBeClearText() bool {
	return true
}

func (r Rendezvous) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(r, m)
	if err != nil {
		return nil, err
	}
	m[rendezvousIdKey] = r.rendezvousId
	return json.Marshal(&m)
}

func DecodeRendezvous(payload []byte) (Rendezvous, error) {
	notiification, err := DecodeRequest(payload)
	if err != nil {
		return Rendezvous{}, err
	}
	rendezvousId, err := jsonparser.GetString(payload, rendezvousIdKey)
	if err != nil {
		return Rendezvous{}, err
	}
	return Rendezvous{notiification, rendezvousId}, nil
}

var _ MoneysocketNotification = &Rendezvous{}
