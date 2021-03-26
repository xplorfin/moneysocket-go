package request

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type OpinionInvoice struct {
	BaseMoneySocketRequest
	ItemID string
}

func NewRequestOpinionInvoice(itemID, requestUUID string) OpinionInvoice {
	r := OpinionInvoice{
		BaseMoneySocketRequest: NewBaseMoneySocketRequest(base.RequestOpinionInvoice),
		ItemID:                 itemID,
	}
	r.BaseUUID = requestUUID
	return r
}

const itemIDKey = "item_id"

func (r OpinionInvoice) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m[itemIDKey] = r.ItemID
	err := EncodeMoneysocketRequest(r, m)
	if err != nil {
		return nil, err
	}
	encodedRequest, err := json.Marshal(m)
	return encodedRequest, err
}

func DecodeRequestOpinionInvoice(payload []byte) (r OpinionInvoice, err error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return OpinionInvoice{}, err
	}

	itemID, err := jsonparser.GetString(payload, itemIDKey)
	if err != nil {
		return OpinionInvoice{}, err
	}
	r = OpinionInvoice{request, itemID}
	return r, nil
}

var _ MoneysocketRequest = &OpinionInvoice{}
