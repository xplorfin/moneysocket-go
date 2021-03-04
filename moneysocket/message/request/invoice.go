package request

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RequestInvoice struct {
	BaseMoneySocketRequest
	Msats int64
}

const msatsKey = "msats"

func NewRequestInvoice(msats int64) RequestInvoice {
	return RequestInvoice{
		NewBaseMoneySocketRequest(base.InvoiceRequest),
		msats,
	}
}

func (r RequestInvoice) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	m[msatsKey] = r.Msats
	err := EncodeMoneysocketRequest(r, m)
	if err != nil {
		return nil, err
	}
	encodedRequest, err := json.Marshal(m)
	return encodedRequest, err
}

// turn a byte slice into a request invoice, return an error if not possible
func DecodeRequestInvoice(payload []byte) (r RequestInvoice, err error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return RequestInvoice{}, err
	}

	msats, err := jsonparser.GetInt(payload, msatsKey)
	if err != nil {
		return RequestInvoice{}, err
	}
	r = RequestInvoice{request, msats}
	return r, nil
}

var _ MoneysocketRequest = &RequestInvoice{}
