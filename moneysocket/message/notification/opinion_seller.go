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
	sellerUUID string
	items      []Item
}

func NewNotifyOpinionSeller(sellerUUID string, items []Item, requestReferenceUUID string) NotifyOpinionSeller {
	return NotifyOpinionSeller{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(msg.NotifyOpinionSeller, requestReferenceUUID),
		sellerUUID:                  sellerUUID,
		items:                       items,
	}
}

const (
	sellerUUIDKey = "seller_uuid"
	itemsKey      = "items"
)

func (o NotifyOpinionSeller) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(o, m)
	if err != nil {
		return nil, err
	}
	m[sellerUUIDKey] = o.sellerUUID
	m[itemsKey] = o.items
	return json.Marshal(&m)
}

func DecodeNotifyOpinionSeller(payload []byte) (NotifyOpinionSeller, error) {
	notification, err := DecodeRequest(payload)
	if err != nil {
		return NotifyOpinionSeller{}, err
	}
	sellerUUID, err := jsonparser.GetString(payload, sellerUUIDKey)
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
		sellerUUID:                  sellerUUID,
		items:                       items,
	}, nil
}

var _ MoneysocketNotification = &NotifyOpinionSeller{}
