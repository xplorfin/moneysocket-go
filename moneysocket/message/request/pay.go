package request

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// Pay is a payment request.
type Pay struct {
	BaseMoneySocketRequest
	Bolt11 string
}

// NewRequestPay creates a new request pay with a given bolt 11
// bolt 11 is not validated client (moneysocket) side.
func NewRequestPay(bolt11 string) Pay {
	return Pay{
		NewBaseMoneySocketRequest(base.PayRequest),
		bolt11,
	}
}

const bolt11key = "Bolt11"

// ToJSON encodes a Pay into json.
func (r Pay) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m[bolt11key] = r.Bolt11
	err := EncodeMoneysocketRequest(r, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// DecodeRequestPay decodes a Pay from json.
func DecodeRequestPay(payload []byte) (r Pay, err error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return Pay{}, err
	}
	bolt11, err := jsonparser.GetString(payload, bolt11key)
	if err != nil {
		return Pay{}, err
	}
	// TODO validate bolt 11 here
	return Pay{
		BaseMoneySocketRequest: request,
		Bolt11:                 bolt11,
	}, nil
}

var _ MoneysocketRequest = &Pay{}
