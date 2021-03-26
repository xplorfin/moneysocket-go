package request

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type Invoice struct {
	BaseMoneySocketRequest
	Msats int64
}

const msatsKey = "msats"

func NewRequestInvoice(msats int64) Invoice {
	return Invoice{
		NewBaseMoneySocketRequest(base.InvoiceRequest),
		msats,
	}
}

func (r Invoice) ToJSON() ([]byte, error) {
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
func DecodeRequestInvoice(payload []byte) (r Invoice, err error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return Invoice{}, err
	}

	msats, err := jsonparser.GetInt(payload, msatsKey)
	if err != nil {
		return Invoice{}, err
	}
	r = Invoice{request, msats}
	return r, nil
}

var _ MoneysocketRequest = &Invoice{}
