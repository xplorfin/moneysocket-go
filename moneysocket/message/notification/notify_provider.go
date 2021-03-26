package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/wad"
)

type NotifyProvider struct {
	BaseMoneySocketNotification
	// account uuid
	AccountUUID string
	// provider pay outgoing invoice
	Payer bool
	// provider generates invoices for incoming payments
	Payee bool
	// balance to advertize as being available
	Wad wad.Wad
}

func NewNotifyProvider(accountUUID string, payer, payee bool, wad wad.Wad, requestUUID string) NotifyProvider {
	return NotifyProvider{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyProvider, requestUUID),
		AccountUUID:                 accountUUID,
		Payer:                       payer,
		Payee:                       payee,
		Wad:                         wad,
	}
}

func (n NotifyProvider) MustBeClearText() bool {
	return false
}

const (
	accountUUIDKey = "account_uuid"
	payerKey       = "payer"
	payeeKey       = "payee"
	wadKey         = "wad"
)

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
