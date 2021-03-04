package notification

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	msg "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// item from the seller
type Item struct {
	// item id
	ItemID string `json:"item_id"`
	// name of the item
	Name string `json:"name"`
	// msats in the message
	Msats int `json:"msats"`
}

type NotifyOpinionSeller struct {
	BaseMoneySocketNotification
	sellerUuid string
	items      []Item
}

func NewNotifyOpinionSeller(sellerUuid string, items []Item, requestReferenceUuid string) NotifyOpinionSeller {
	return NotifyOpinionSeller{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(msg.NotifyOpinionSeller, requestReferenceUuid),
		sellerUuid:                  sellerUuid,
		items:                       items,
	}
}

const (
	sellerUuidKey = "seller_uuid"
	itemsKey      = "items"
)

func (o NotifyOpinionSeller) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(o, m)
	if err != nil {
		return nil, err
	}
	m[sellerUuidKey] = o.sellerUuid
	m[itemsKey] = o.items
	return json.Marshal(&m)
}

func DecodeNotifyOpinionSeller(payload []byte) (NotifyOpinionSeller, error) {
	notification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyOpinionSeller{}, err
	}
	sellerUuid, err := jsonparser.GetString(payload, sellerUuidKey)
	if err != nil {
		return NotifyOpinionSeller{}, err
	}
	rawItems, _, _, err := jsonparser.Get(payload, itemsKey)
	if err != nil {
		return NotifyOpinionSeller{}, err
	}
	var items []Item
	err = json.Unmarshal(rawItems, &items)
	if err != nil {
		return NotifyOpinionSeller{}, err
	}
	return NotifyOpinionSeller{
		BaseMoneySocketNotification: notification,
		sellerUuid:                  sellerUuid,
		items:                       items,
	}, nil
}

var _ MoneysocketNotification = &NotifyOpinionSeller{}
