package notification

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type NotifyPreimage struct {
	BaseMoneySocketNotification
	Preimage string
	Ext      string
}

func NewNotifyPreimage(preimage, ext, requestUuid string) NotifyPreimage {
	return NotifyPreimage{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyPreimage, requestUuid),
		Preimage:                    preimage,
		Ext:                         ext,
	}
}

func (n NotifyPreimage) MustBeClearText() bool {
	return false
}

const (
	preimageKey = "preimage"
	extKey      = "ext"
)

func (n NotifyPreimage) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	m[preimageKey] = n.Preimage
	m[extKey] = n.Ext
	return json.Marshal(&m)
}

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
