package request

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RequestOpinionInvoice struct {
	BaseMoneySocketRequest
	ItemId string
}

func NewRequestOpinionInvoice(itemId, requestUuid string) RequestOpinionInvoice {
	r := RequestOpinionInvoice{
		BaseMoneySocketRequest: NewBaseMoneySocketRequest(base.RequestOpinionInvoice),
		ItemId:                 itemId,
	}
	r.BaseUuid = requestUuid
	return r
}

const itemIdKey = "item_id"

func (r RequestOpinionInvoice) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	m[itemIdKey] = r.ItemId
	err := EncodeMoneysocketRequest(r, m)
	if err != nil {
		return nil, err
	}
	encodedRequest, err := json.Marshal(m)
	return encodedRequest, err
}

func DecodeRequestOpinionInvoice(payload []byte) (r RequestOpinionInvoice, err error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return RequestOpinionInvoice{}, err
	}

	itemId, err := jsonparser.GetString(payload, itemIdKey)
	if err != nil {
		return RequestOpinionInvoice{}, err
	}
	r = RequestOpinionInvoice{request, itemId}
	return r, nil
}

var _ MoneysocketRequest = &RequestOpinionInvoice{}
