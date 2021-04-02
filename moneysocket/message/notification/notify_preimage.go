package notification

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// NotifyPreimage is a notification that a preimage is ready
type NotifyPreimage struct {
	BaseMoneySocketNotification
	Preimage string
	Ext      string
}

// NewNotifyPreimage creates a notify preimages
func NewNotifyPreimage(preimage, ext, requestUUID string) NotifyPreimage {
	return NotifyPreimage{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyPreimage, requestUUID),
		Preimage:                    preimage,
		Ext:                         ext,
	}
}

// MustBeClearText determines wether the message must be clear text
func (n NotifyPreimage) MustBeClearText() bool {
	return false
}

const (
	preimageKey = "preimage"
	extKey      = "ext"
)

// ToJSON converts the message to json
func (n NotifyPreimage) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	m[preimageKey] = n.Preimage
	m[extKey] = n.Ext
	return json.Marshal(&m)
}

// IsValid determines whether the notification is valid
func (n NotifyPreimage) IsValid() (bool, error) {
	_, err := strconv.ParseUint(n.Preimage, 16, 64)
	if err != nil {
		return false, fmt.Errorf("preimage must be a hex string")
	}
	if len(n.Preimage) != 64 {
		return false, fmt.Errorf("preimage not 256-bit value hex string")
	}
	return true, err
}

var _ MoneysocketNotification = &NotifyPreimage{}
