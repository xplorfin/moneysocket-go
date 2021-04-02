package message

import (
	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
)

// FromText generates a moneysocket message from a layer
func FromText(payload []byte) (base.MoneysocketMessage, base.MessageType, error) {
	class, err := jsonparser.GetString(payload, MessageClass)
	if err != nil {
		return nil, 0, err
	}
	if class != base.Notification.ToString() {
		return request.FromText(payload)
	}
	panic("method not yet implemented")
}
