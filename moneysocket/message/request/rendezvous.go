package request

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// Rendezvous is a rendezvous request.
type Rendezvous struct {
	BaseMoneySocketRequest
	// RendezvousID is the id of the rendezvous we're requesting (normally derived  from the shared seed)
	RendezvousID string
}

// NewRendezvousRequest is the request the server start a rendezvous w/ a given rendezvous id.
func NewRendezvousRequest(id string) Rendezvous {
	return Rendezvous{
		NewBaseMoneySocketRequest(base.RendezvousRequest),
		id,
	}
}

const rendezvousIDKey = "rendezvous_id"

// ToJSON marshals the Rendezvous message.
func (r Rendezvous) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m[rendezvousIDKey] = r.RendezvousID
	err := EncodeMoneysocketRequest(r, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// MustBeClearText determines that a rendezvous message can be a clear text message.
func (r Rendezvous) MustBeClearText() bool {
	return true
}

// DecodeRendezvousRequest decodes a Rendezvous request from a payload.
func DecodeRendezvousRequest(payload []byte) (r Rendezvous, err error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return Rendezvous{}, err
	}

	rid, err := jsonparser.GetString(payload, rendezvousIDKey)
	if err != nil {
		return Rendezvous{}, nil
	}
	r = Rendezvous{
		BaseMoneySocketRequest: request,
		RendezvousID:           rid,
	}
	return r, nil
}

var _ MoneysocketRequest = &Rendezvous{}
