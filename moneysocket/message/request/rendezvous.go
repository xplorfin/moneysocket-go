package request

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RequestRendezvous struct {
	BaseMoneySocketRequest
	// the id of the rendezvous we're requesting (normally derived  from the shared seed)
	RendezvousId string
}

// request the server start a rendezvous w/ a given rendezvous id
func NewRendezvousRequest(id string) RequestRendezvous {
	return RequestRendezvous{
		NewBaseMoneySocketRequest(base.RendezvousRequest),
		id,
	}
}

const rendevousIdKey = "rendezvous_id"

func (r RequestRendezvous) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	m[rendevousIdKey] = r.RendezvousId
	err := EncodeMoneysocketRequest(r, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (r RequestRendezvous) MustBeClearText() bool {
	return true
}

func DecodeRendezvousRequest(payload []byte) (r RequestRendezvous, err error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return RequestRendezvous{}, err
	}

	rid, err := jsonparser.GetString(payload, rendevousIdKey)
	if err != nil {
		return RequestRendezvous{}, nil
	}
	r = RequestRendezvous{
		BaseMoneySocketRequest: request,
		RendezvousId:           rid,
	}
	return r, nil
}

var _ MoneysocketRequest = &RequestRendezvous{}
