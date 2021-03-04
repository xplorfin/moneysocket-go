package notification

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/wad"
)

type NotifyProvider struct {
	BaseMoneySocketNotification
	// account uuid
	AccountUuid string
	// provider pay outgoing invoice
	Payer bool
	// provider generates invoices for incoming payments
	Payee bool
	// balance to advertize as being available
	Wad wad.Wad
}

func NewNotifyProvider(accountUuid string, payer, payee bool, wad wad.Wad, requestUuid string) NotifyProvider {
	return NotifyProvider{
		BaseMoneySocketNotification: NewBaseMoneySocketNotification(base.NotifyProvider, requestUuid),
		AccountUuid:                 accountUuid,
		Payer:                       payer,
		Payee:                       payee,
		Wad:                         wad,
	}
}

func (n NotifyProvider) MustBeClearText() bool {
	return false
}

const (
	accountUuidKey = "account_uuid"
	payerKey       = "payer"
	payeeKey       = "payee"
	wadKey         = "wad"
)

func (n NotifyProvider) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketNotification(n, m)
	if err != nil {
		return nil, err
	}
	m[accountUuidKey] = n.AccountUuid
	m[payerKey] = n.Payer
	m[payeeKey] = n.Payee
	m[wadKey] = n.Wad
	return json.Marshal(&m)
}
