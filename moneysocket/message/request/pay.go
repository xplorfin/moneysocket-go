package request

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RequestPay struct {
	BaseMoneySocketRequest
	Bolt11 string
}

// create a new request pay with a given bolt 11
// bolt 11 is not validated client (moneysocket) side
func NewRequestPay(bolt11 string) RequestPay {
	return RequestPay{
		NewBaseMoneySocketRequest(base.PayRequest),
		bolt11,
	}
}

const bolt11key = "Bolt11"

func (r RequestPay) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	m[bolt11key] = r.Bolt11
	err := EncodeMoneysocketRequest(r, m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func DecodeRequestPay(payload []byte) (r RequestPay, err error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return RequestPay{}, err
	}
	bolt11, err := jsonparser.GetString(payload, bolt11key)
	if err != nil {
		return RequestPay{}, err
	}
	// TODO validate bolt 11 here
	return RequestPay{
		BaseMoneySocketRequest: request,
		Bolt11:                 bolt11,
	}, nil
}

var _ MoneysocketRequest = &RequestPay{}
