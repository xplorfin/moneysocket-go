package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/wad"
)

// NotifyProvider notifies a provider is ready
type NotifyProvider struct {
	BaseMoneySocketNotification
	// AccountUUID
	AccountUUID string
	// Payer is the provider pay outgoing invoice
	Payer bool
	// Payee is the provider generates invoices for incoming payments
	Payee bool
	// Wad is the balance to advertise as being available
	Wad wad.Wad
}

// NewNotifyProvider creates a NotifyProvider notification
func NewNotifyProvider(accountUUID string, payer, payee bool, wad wad.Wad, requestUUID string) NotifyProvider {
	return NotifyProvider{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyProvider, requestUUID),
		AccountUUID:                 accountUUID,
		Payer:                       payer,
		Payee:                       payee,
		Wad:                         wad,
	}
}

// MustBeClearText is always false  for NotifyProvider messages
func (n NotifyProvider) MustBeClearText() bool {
	return false
}

const (
	accountUUIDKey = "account_uuid"
	payerKey       = "payer"
	payeeKey       = "payee"
	wadKey         = "wad"
)

// ToJSON converts a message to a json payload
func (n NotifyProvider) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneySocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	m[accountUUIDKey] = n.AccountUUID
	m[payerKey] = n.Payer
	m[payeeKey] = n.Payee
	m[wadKey] = n.Wad
	return json.Marshal(&m)
}
